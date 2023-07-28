package sensor

import (
	"context"
	"database/sql"
	"fmt"
	"sensors-generator/config"
	"sensors-generator/internal/apperror"
	clients "sensors-generator/pkg/client/interfaces"
	"sensors-generator/pkg/logging"
	"time"
)

type repository struct {
	client clients.DBClient
	logger *logging.Logger
	cfg    *config.Config
}

func NewPostgresqlRepository(client *sql.DB,
	logger *logging.Logger, cfg *config.Config) *repository {
	return &repository{
		client: client,
		logger: logger,
		cfg:    cfg,
	}
}

func (r *repository) FindAll(ctx context.Context, filters SensorFilters) ([]Sensor, error) {
	q := `SELECT s.id, sg.name, s.index, s.x, s.y, s.z, s.data_output_rate, s.created_at, s.updated_at FROM sensors as s
		JOIN sensor_groups sg ON s.group_id=sg.id`

	rows, err := r.client.QueryContext(ctx, q)
	defer rows.Close()
	if err != nil {
		r.logger.Errorf("Failed to get rows, due to error: %v", err)
		return nil, apperror.ErrorWithMessage(apperror.ErrBadRequest, "Sensors not found.")
	}

	sensors := make([]Sensor, 0)

	for rows.Next() {
		var sensor Sensor
		if err := rows.Scan(&sensor.ID, &sensor.CodeName.GroupName, &sensor.CodeName.Index, &sensor.Coords.X,
			&sensor.Coords.Y, &sensor.Coords.Z, &sensor.DataOutputRate, &sensor.CreatedAt, &sensor.UpdatedAt); err != nil {
			r.logger.Errorf("Failed to fetch row, due to error: %v", err)
			return nil, apperror.ErrInternalSystem
		}

		sensors = append(sensors, sensor)
	}

	return sensors, nil
}

func (r *repository) FindOneByID(ctx context.Context, id int, filters SensorFilters) (*Sensor, error) {
	return nil, nil
}

func (r *repository) Create(ctx context.Context, sensor CreateSensorDTO) error {
	query := "SELECT id FROM sensor_groups WHERE name = $1 LIMIT 1"
	var groupID int
	err := r.client.QueryRowContext(ctx, query, sensor.CodeName.GroupName).Scan(&groupID)
	if err != nil {
		r.logger.Errorf("Failed to retrieve group ID, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	insertQuery := `INSERT INTO sensors (group_id, index, x, y, z, data_output_rate, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`

	t := time.Now()

	if _, err := r.client.ExecContext(ctx, insertQuery, groupID, sensor.CodeName.Index,
		sensor.Coords.X, sensor.Coords.Y, sensor.Coords.Z, sensor.DataOutputRate,
		t, t); err != nil {
		r.logger.Errorf("Failed create sensor, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	return nil
}

func (r *repository) AddSensorToGroup(ctx context.Context, sensorID int, groupID int) error {
	q := `UPDATE sensors SET group_id=$1 WHERE id=$2`

	if _, err := r.client.ExecContext(ctx, q, groupID, sensorID); err != nil {
		r.logger.Errorf("Failed to add sensor to another group, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	return nil
}

func (r *repository) FindMaxTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates) (float32, error) {
	q := `SELECT MAX(sd.temperature) FROM sensors as sens
		JOIN sensor_data sd ON sens.id=sd.sensor_id
		WHERE sens.x < $1 AND sens.x > $2 AND sens.y < $3 AND sens.y > $4 AND sens.z < $5 AND sens.z > $6`

	var temperature float32

	if err := r.client.QueryRow(q, maxCoords.X, minCoords.X,
		maxCoords.Y, minCoords.Y, maxCoords.Z, minCoords.Z).Scan(&temperature); err != nil {
		r.logger.Errorf("Cannot measure max temperature, due to error: %v", err)
		return 0.0, apperror.ErrorWithMessage(apperror.ErrInternalSystem, "No data was found.")
	}

	return temperature, nil
}

func (r *repository) FindMinTemperatureForRegion(ctx context.Context, minCoords, maxCoords Coordinates) (float32, error) {
	q := `SELECT MIN(sd.temperature) FROM sensors as sens
		JOIN sensor_data sd ON sens.id=sd.sensor_id
		WHERE sens.x < $1 AND sens.x > $2 AND sens.y < $3 AND sens.y > $4 AND sens.z < $5 AND sens.z > $6`

	var temperature float32

	if err := r.client.QueryRow(q, maxCoords.X, minCoords.X,
		maxCoords.Y, minCoords.Y, maxCoords.Z, minCoords.Z).Scan(&temperature); err != nil {
		r.logger.Errorf("Cannot measure min temperature, due to error: %v", err)
		return 0.0, apperror.ErrorWithMessage(apperror.ErrInternalSystem, "No data was found.")
	}

	return temperature, nil
}

func (r *repository) FindAvgTemperatureForSensor(ctx context.Context, filters SensorFilters) (float32, error) {
	q := `SELECT AVG(sd.temperature) FROM sensors AS sens
		JOIN sensor_data sd ON sens.id=sd.sensor_id`

	args := []interface{}{}
	argsCounter := 1

	if !filters.CodeName.IsEmpty() {
		q += "\n" + `JOIN sensor_groups sg ON sg.id=sens.group_id`
		if argsCounter == 1 {
			q += "\n" + ` WHERE`
		} else {
			q += ` AND`
		}
		q += fmt.Sprintf(` sg.name=$%d AND sens.index=$%d`, argsCounter, argsCounter+1)
		args = append(args, filters.CodeName.GroupName)
		args = append(args, filters.CodeName.Index)
		argsCounter += 2
	}

	if !filters.FromDate.IsZero() {
		if argsCounter == 1 {
			q += "\n" + ` WHERE`
		} else {
			q += ` AND`
		}
		q += fmt.Sprintf(` sd.created_at >= $%d`, argsCounter)
		args = append(args, filters.FromDate)
		argsCounter++
	}

	if !filters.TillDate.IsZero() {
		if argsCounter == 1 {
			q += "\n" + ` WHERE`
		} else {
			q += ` AND`
		}
		q += fmt.Sprintf(` sd.created_at <= $%d`, argsCounter)
		args = append(args, filters.TillDate)
		argsCounter++
	}

	// if !filters.FromDate.IsZero() && !filters.TillDate.IsZero() {
	// 	if argsCounter == 1 {
	// 		q += "\n" + ` WHERE`
	// 	} else {
	// 		q += ` AND`
	// 	}
	// 	q += fmt.Sprintf(` sd.created_at BETWEEN $%d AND $%d`, argsCounter, argsCounter+1)
	// 	args = append(args, filters.FromDate)
	// 	args = append(args, filters.TillDate)
	// 	argsCounter += 2
	// }

	var temperature float32

	if err := r.client.QueryRowContext(ctx, q, args...).Scan(&temperature); err != nil {
		r.logger.Errorf("Cannot measure average temperature, due to error: %v", err)
		return 0.0, apperror.ErrorWithMessage(apperror.ErrInternalSystem, "No sensor data was found.")
	}

	return temperature, nil
}
