package database

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/kosdirus/andintern/internal/config"

	_ "github.com/jackc/pgx/v4/stdlib"
)

type Client struct {
	*sqlx.DB

	schemaName string
}

func NewClient(cfg config.Config) (*Client, error) {
	db, err := sqlx.Open("pgx", cfg.DB.URL)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.DB.MaxOpenConns)
	db.SetMaxIdleConns(cfg.DB.MaxIdleConns)

	return &Client{
		DB:         db,
		schemaName: cfg.DB.SchemaName,
	}, nil
}

func (db *Client) Migrate() error {
	if _, err := db.Exec(`CREATE SCHEMA IF NOT EXISTS ` + db.schemaName); err != nil {
		return fmt.Errorf("can't create schema: %w", err)
	}

	m, err := migrations(db.schemaName, "migrations")
	if err != nil {
		return fmt.Errorf("can't create a new migrator instance: %w", err)
	}

	// Migrate up
	if err := m.Migrate(db.DB.DB); err != nil {
		return fmt.Errorf("can't migrate the DB: %w", err)
	}

	return nil
}
