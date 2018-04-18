package sql

import (
	"database/sql"

	"github.com/gobuffalo/packr"
	migrate "github.com/rubenv/sql-migrate"
)

func maybeMigrate(db *sql.DB, dbType string) error {
	migrations := &migrate.PackrMigrationSource{
		Box: packr.NewBox("./migrations"),
	}

	_, err := migrate.Exec(db, dbType, migrations, migrate.Up)
	if err != nil {
		return err
	}

	return nil
}
