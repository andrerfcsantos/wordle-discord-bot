package db

import (
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"os"
)

func (r *Repository) RunMigrations() error {

	var migrationPath = "db/migrations"
	switch os.Getenv("WORDLE_ENVIRONMENT") {
	case "dev":
		migrationPath = "db/migrations/dev"
	}

	err := r.runMigrations(migrationPath)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) runMigrations(migrationPath string) error {

	database, err := r.db.DB()
	if err != nil {
		return fmt.Errorf("getting sql.DB for migrations: %v", err)
	}

	driver, err := postgres.WithInstance(database, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("creating migrate driver: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationPath),
		"postgres",
		driver,
	)

	if err != nil {
		return fmt.Errorf("creating migrate: %v", err)
	}

	err = m.Up()
	if err != nil {
		return fmt.Errorf("running migrations: %v", err)
	}

	return nil
}
