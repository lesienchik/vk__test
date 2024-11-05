package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/lesienchik/vk__test/internal/config"
)

func ConnectToDb(cfg *config.Postgres) (*sql.DB, error) {
	DSN := fmt.Sprintf(
		"dbname=%s user=%s password=%s host=%s port=%d sslmode=%s fallback_application_name=%s",
		cfg.DbName,
		cfg.User,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Sslmode,
		cfg.AppName,
	)

	db, err := sql.Open("postgres", DSN)
	if err != nil {
		return nil, fmt.Errorf("ConnectToDb (1): %w", err)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ConnectToDb (2): %w", err)
	}

	db.SetMaxOpenConns(cfg.MaxConns)

	return db, nil
}
