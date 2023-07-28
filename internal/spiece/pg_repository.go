package spiece

import (
	"context"
	"database/sql"
	"sensors-generator/config"
	"sensors-generator/internal/apperror"
	clients "sensors-generator/pkg/client/interfaces"
	"sensors-generator/pkg/logging"
	"strconv"
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

func (r *repository) FindAll(ctx context.Context, filters SpieceFilters) ([]Spiece, error) {
	q := ``
	args := make([]interface{}, 0)
	argsCounter := 1

	if len(filters.GroupName) > 0 {
		q += `SELECT s.id, s.name, s.created_at, s.updated_at FROM sensors as sens
			JOIN sensor_data sd ON sens.id=sd.sensor_id
			JOIN detected_spieces ds ON sd.id=ds.sensor_data_id
			JOIN spieces s ON s.id=ds.spiece_id
			WHERE sens.group_name=$` + strconv.Itoa(argsCounter)

		args = append(args, filters.GroupName)
		argsCounter++
	} else {
		q += `SELECT id, name, created_at, updated_at FROM spieces`
	}

	if filters.N > 0 {
		q += ` ORDER BY created_at DESC LIMIT $` + strconv.Itoa(argsCounter)
		args = append(args, filters.N)
		argsCounter++
	}

	var spieces []Spiece

	rows, err := r.client.Query(q, args...)
	defer rows.Close()
	if err != nil {
		r.logger.Errorf("Cannot find spieces, due to error: %v", err)
		return nil, apperror.ErrorWithMessage(apperror.ErrBadRequest, "Spieces not found.")
	}

	for rows.Next() {
		var spiece Spiece
		if err := rows.Scan(&spiece.ID, &spiece.Name,
			&spiece.CreatedAt, &spiece.UpdatedAt); err != nil {
			r.logger.Errorf("Cannot scan spieces row.")
			return nil, apperror.ErrInternalSystem
		}
		spieces = append(spieces, spiece)
	}

	return spieces, nil
}

func (r *repository) FindOneByID(ctx context.Context, id int, filters SpieceFilters) (*Spiece, error) {
	q := `SELECT id, name, created_at, updated_at FROM spieces WHERE id=$1`
	var spiece Spiece

	if err := r.client.QueryRow(q, id).Scan(&spiece.ID, &spiece.Name,
		&spiece.CreatedAt, &spiece.UpdatedAt); err != nil {
		r.logger.Errorf("Cannot scan spieces row.")
		return nil, apperror.ErrorWithMessage(apperror.ErrBadRequest, "Spiece not found.")
	}

	return &spiece, nil
}

func (r *repository) Create(ctx context.Context, spiece CreateSpieceDTO) error {
	q := `INSERT INTO spieces(name, created_at, updated_at)
		VALUES($1, $2, $3)`

	t := time.Now()

	if _, err := r.client.ExecContext(ctx, q, spiece.Name, t, t); err != nil {
		r.logger.Errorf("Cannot create spiece, due to error: %v", err)
		return apperror.ErrInternalSystem
	}

	return nil
}
