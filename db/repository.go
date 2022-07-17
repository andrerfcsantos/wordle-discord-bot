package db

import (
	"fmt"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"os"
	"time"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository() (*Repository, error) {
	var err error
	var repo Repository

	repo.db, err = getConnection()
	if err != nil {
		return nil, fmt.Errorf("getting database session: %w", err)
	}

	return &repo, nil
}

func getConnection() (*gorm.DB, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("WORDLE_DB_HOST"),
		os.Getenv("WORDLE_DB_USER"),
		os.Getenv("WORDLE_DB_PASSWORD"),
		os.Getenv("WORDLE_DB_NAME"),
		os.Getenv("WORDLE_DB_PORT"))

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("getting database session: %w", err)
	}

	return db, nil
}

func (r *Repository) Database() *gorm.DB {
	db, err := r.db.DB()
	if err != nil {
		log.Errorf("getting database session in Database(): %v", err)
	}

	errPing := db.Ping()
	for errPing != nil {
		log.Errorf("could not ping database in Database(): %v", err)
		time.Sleep(time.Second * 5)
		errPing = db.Ping()
	}
	return r.db
}
