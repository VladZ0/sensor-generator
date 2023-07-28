package group

import (
	"context"
	"database/sql"
	"fmt"
	"sensors-generator/config"
	"sensors-generator/internal/apperror"
	"sensors-generator/internal/spiece"
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

func (r *repository) FindSpiecesInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (map[*spiece.Spiece]int, error) {
	q := ``

	detectedSpiecesTbl := "detected_spieces"

	argsCounter := 1
	args := []interface{}{}

	if filters.TopLimit > 0 {
		ds := fmt.Sprintf(`WITH det_s AS (
			SELECT ds.spiece_id, ds.sensor_data_id
			FROM detected_spieces ds
			JOIN sensor_data sd ON ds.sensor_data_id=sd.id
			JOIN sensors sens ON sd.sensor_id=sens.id
			JOIN sensor_groups sg ON sg.id=sens.group_id
			WHERE sg.name='%v'
			ORDER BY sd.created_at DESC
			LIMIT %d
		)`, groupName, filters.TopLimit)
		q += ds
		detectedSpiecesTbl = "det_s"
	}

	q += fmt.Sprintf(`SELECT s.id, s.name, s.created_at, s.updated_at, COUNT(s.id)
	FROM sensor_groups as sg
	JOIN sensors sens ON sg.id = sens.group_id
	JOIN sensor_data sd ON sens.id = sd.sensor_id
	JOIN %s ds ON sd.id = ds.sensor_data_id
	JOIN spieces s ON s.id = ds.spiece_id
	WHERE sg.name = $%d`, detectedSpiecesTbl, argsCounter)

	args = append(args, groupName)
	argsCounter++

	if !filters.FromDate.IsZero() && !filters.TillDate.IsZero() {
		q += fmt.Sprintf(` AND sd.created_at BETWEEN $%d AND $%d`, argsCounter, argsCounter+1)
		args = append(args, filters.FromDate)
		args = append(args, filters.TillDate)
		argsCounter += 2
	}

	q += "\n" + `GROUP BY s.id, s.name, s.created_at, s.updated_at`

	spieces := make(map[*spiece.Spiece]int)

	rows, err := r.client.Query(q, args...)
	defer rows.Close()
	if err != nil {
		r.logger.Errorf("Cannot find spieces, due to error: %v", err)
		return nil, apperror.ErrInternalSystem
	}

	var count int

	for rows.Next() {
		var spiece spiece.Spiece
		if err := rows.Scan(&spiece.ID, &spiece.Name,
			&spiece.CreatedAt, &spiece.UpdatedAt, &count); err != nil {
			r.logger.Errorf("Cannot scan spieces row.")
			return nil, apperror.ErrInternalSystem
		}
		spieces[&spiece] = count
	}

	return spieces, nil
}

func (r *repository) FindOneByID(ctx context.Context, id int, filters SensorGroupFilters) (*SensorGroup, error) {
	q := `SELECT id, name, created_at, updated_at FROM sensor_groups WHERE id=$1`
	var sensorGroup SensorGroup

	if err := r.client.QueryRow(q, id).Scan(&sensorGroup.ID, &sensorGroup.Name,
		&sensorGroup.CreatedAt, &sensorGroup.UpdatedAt); err != nil {
		r.logger.Errorf("Cannot scan sensor group row.")
		return nil, apperror.ErrInternalSystem
	}

	return &sensorGroup, nil
}

func (r *repository) FindAll(ctx context.Context, filters SensorGroupFilters) ([]SensorGroup, error) {
	q := `SELECT id, name, created_at, updated_at FROM sensor_groups`

	rows, err := r.client.QueryContext(ctx, q)
	if err != nil {
		r.logger.Errorf("Failed to get rows, due to error: %v", err)
		return nil, apperror.ErrInternalSystem
	}
	defer rows.Close()

	sensorGroups := make([]SensorGroup, 0)

	for rows.Next() {
		var sensorGroup SensorGroup
		if err := rows.Scan(&sensorGroup.ID, &sensorGroup.Name, &sensorGroup.CreatedAt, &sensorGroup.UpdatedAt); err != nil {
			r.logger.Errorf("Failed to fetch row, due to error: %v", err)
			return nil, apperror.ErrInternalSystem
		}

		sensorGroups = append(sensorGroups, sensorGroup)
	}

	return sensorGroups, nil
}

func (r *repository) FindAvgTransparencyInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (uint8, error) {
	q := `SELECT AVG(sd.transparency) FROM sensor_groups as sg
		JOIN sensors sens ON sg.id=sens.group_id
		JOIN sensor_data sd ON sens.id=sd.sensor_id
		WHERE sg.name=$1`

	args := []interface{}{
		groupName,
	}
	// argsCounter := 2
	var transparency float32

	if err := r.client.QueryRow(q, args...).Scan(&transparency); err != nil {
		r.logger.Errorf("Cannot measure transparency, due to error: %v", err)
		return 0, apperror.ErrInternalSystem
	}

	r.logger.Info("Transparency were successfully measured.")
	return uint8(transparency), nil
}

func (r *repository) FindAvgTemperatureInGroup(ctx context.Context, groupName string, filters SensorGroupFilters) (float32, error) {
	r.logger.Info("FIND AVERAGE TEMPERATURE IN GROUP.")
	q := `SELECT AVG(sd.temperature) FROM sensor_groups as sg
		JOIN sensors sens ON sg.id=sens.group_id
		JOIN sensor_data sd ON sens.id=sd.sensor_id
		WHERE sg.name=$1`

	args := []interface{}{
		groupName,
	}
	// argsCounter := 2
	var temperature float32

	if err := r.client.QueryRow(q, args...).Scan(&temperature); err != nil {
		r.logger.Errorf("Cannot measure transparency, due to error: %v", err)
		return 0.0, apperror.ErrInternalSystem
	}

	return temperature, nil
}

func (r *repository) Create(ctx context.Context, grp CreateSensorGroupDTO) error {
	q := `INSERT INTO sensor_groups(name, created_at, updated_at)
			VALUES($1, $2, $3)`

	t := time.Now()

	if _, err := r.client.ExecContext(ctx, q, grp.Name, t, t); err != nil {
		r.logger.Errorf("Cannot create group, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	r.logger.Info("Group was created successfully.")
	return nil
}
