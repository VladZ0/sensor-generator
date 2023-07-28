package sensordata

import (
	"context"
	"database/sql"
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

func (r *repository) FindAll(ctx context.Context, filters SensorDataFilters) ([]SensorData, error) {
	return nil, nil
}

func (r *repository) FindOneByID(ctx context.Context, id int, filters SensorDataFilters) (*SensorData, error) {
	q := `SELECT sd.id, sens.id, sd.temperature, sd.transparency, sd.created_at, sd.updated_at FROM sensor_data AS sd
		JOIN sensors sens ON sd.sensor_id=sens.id
		WHERE sd.id=$1`

	qDetectedSpieces := `SELECT s.id, s.name, s.created_at, s.updated_at FROM sensor_data AS sd
	JOIN detected_spieces ds ON sd.id=ds.sensor_data_id
	JOIN spieces s ON ds.spiece_id=s.id
	WHERE sd.id=$1`

	tx, err := r.client.BeginTx(ctx, &sql.TxOptions{})

	var sensorData SensorData

	err = tx.QueryRowContext(ctx, q, id).Scan(&sensorData.ID, &sensorData.SensorID,
		&sensorData.Temperature, &sensorData.Transparency, &sensorData.CreatedAt, &sensorData.UpdatedAt)

	sensorData.DetectedSpieces = make([]spiece.Spiece, 0)

	rows, err := tx.QueryContext(context.Background(), qDetectedSpieces, id)
	defer rows.Close()

	for rows.Next() {
		var detectedSpiece spiece.Spiece
		if err := rows.Scan(&detectedSpiece.ID, &detectedSpiece.Name,
			&detectedSpiece.CreatedAt, &detectedSpiece.UpdatedAt); err != nil {
			tx.Rollback()
			r.logger.Errorf("Cannot find detected spiece, due to error: %v", err)
			return nil, nil
		}
		sensorData.DetectedSpieces = append(sensorData.DetectedSpieces, detectedSpiece)
	}

	switch err {
	case nil:
		tx.Commit()
		break

	default:
		tx.Rollback()
		r.logger.Errorf("Cannot find sensor data, due to error: %v", err)
		return nil, nil
	}

	return &sensorData, nil
}

func (r *repository) Create(ctx context.Context, sensorData CreateSensorDataDTO) (int, error) {
	q := `INSERT INTO sensor_data(sensor_id, temperature, transparency, created_at, updated_at)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id`

	t := time.Now()

	var id int

	if err := r.client.QueryRowContext(ctx, q, sensorData.SensorID, sensorData.Temperature,
		sensorData.Transparency, t, t).Scan(&id); err != nil {
		r.logger.Errorf("Cannot create sensor data, due to error: %v", err)
		return 0, apperror.ErrInternalSystem
	}

	return id, nil
}

func (r *repository) AddDetectedSpiece(ctx context.Context, sensorDataID int, spiece spiece.Spiece) error {
	q := `INSERT INTO detected_spieces(spiece_id, sensor_data_id)
		VALUES($1, $2)`

	if _, err := r.client.ExecContext(ctx, q, spiece.ID, sensorDataID); err != nil {
		r.logger.Errorf("Failed to detect spiece, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	return nil
}
