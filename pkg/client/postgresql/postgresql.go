package postgresql

import (
	"database/sql"
	"fmt"
	"sensors-generator/pkg/logging"

	_ "github.com/lib/pq"
)

type PgConfig struct {
	Username string `yaml:"username" env:"USERNAME" env-default:"vlad"`
	Database string `yaml:"database" env:"DATABASE" env-default:"messenger"`
	Password string `yaml:"password" env:"PASSWORD" env-default:""`
	Host     string `yaml:"host" env:"HOST" env-default:"localhost"`
	Port     string `yaml:"port" env:"PORT" env-default:"5432"`
	SSLMode  string `yaml:"ssl_mode" env:"SSL_MODE" env-default:"disable"`
}

func NewClient(cfg PgConfig) (*sql.DB, error) {
	var dsn string

	dsn = fmt.Sprintf(
		"user=%s dbname=%s password=%s host=%s port=%s  sslmode=%s",
		cfg.Username, cfg.Database, cfg.Password, cfg.Host, cfg.Port, cfg.SSLMode,
	)

	logging.GetLogger().Infof("%v", cfg.Database)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
