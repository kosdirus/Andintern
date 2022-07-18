package database

import (
	"fmt"
	"github.com/lopezator/migrator"
)

func migrations(schema, migrationsTable string) (*migrator.Migrator, error) {
	return migrator.New(migrator.TableName(fmt.Sprintf("%s.%s", schema, migrationsTable)),
		migrator.Migrations(
			migrationInit(schema),
		),
	)
}
