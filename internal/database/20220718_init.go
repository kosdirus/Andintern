package database

import (
	"database/sql"
	"fmt"
	"github.com/lopezator/migrator"
)

//nolint // to bypass gosec sql concat warning
func migrationInit(schema string) *migrator.Migration {
	return &migrator.Migration{
		Name: "20220718_init",
		Func: func(tx *sql.Tx) error {
			query := `CREATE TABLE IF NOT EXISTS ` + schema + `.cars (` +
				`  id SERIAL` +
				`, brand VARCHAR(128) NOT NULL` +
				`, price BIGINT NOT NULL` +
				`, CONSTRAINT cars_pkey PRIMARY KEY (brand)` +
				`)`

			if _, err := tx.Exec(query); err != nil {
				return fmt.Errorf("applying 20220718_init migration: %w", err)
			}

			return nil
		},
	}
}
