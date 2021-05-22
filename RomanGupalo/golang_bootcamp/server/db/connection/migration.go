package connection

import (
	"database/sql"
	"errors"
	"log"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// TryMigrate trying to execute migration scripts from migrations_source
// Returns nil on success
func TryMigrate(migrations_source string) error {
	if session == nil {
		return errors.New("session not created")
	}

	log.Println("Starting migrations")
	db := session.Driver().(*sql.DB)
	driver, err := mysql.WithInstance(db, &mysql.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(migrations_source, "mysql", driver)
	if err != nil {
		return err
	}

	switch err = m.Up(); err {
	case nil:
		log.Println("Migrations executed")
	case migrate.ErrNoChange:
		log.Println("No migrations to execute")
	default:
		return err
	}

	return nil
}
