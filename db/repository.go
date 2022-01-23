package db

import (
	"fmt"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type Attempt struct {
	ChannelId    string    `gorm:"primary_key;column:channel_id"`
	UserId       string    `gorm:"primary_key;column:user_id"`
	Day          int       `gorm:"primary_key;column:day"`
	MessageId    string    `gorm:"column:message_id"`
	UserName     string    `gorm:"column:user_name"`
	Attempts     int       `gorm:"column:attempts"`
	MaxAttempts  int       `gorm:"column:max_attempts"`
	Success      bool      `gorm:"column:success"`
	AttemptsJson string    `gorm:"column:attempts_json"`
	PostedAt     time.Time `gorm:"column:posted_at"`
	Score        int       `gorm:"column:score"`
}

type TrackedChannel struct {
	ChannelId string `gorm:"primary_key;column:channel_id"`
}

type Repository struct {
	db *gorm.DB
}

func (r *Repository) IsTrackedChannel(channelId string) (bool, error) {
	var count int64

	query := r.db.
		Table("tracked_channels").
		Where("channel_id = ?", channelId).
		Count(&count)

	if query.Error != nil {
		return false, fmt.Errorf("checking tracked channel: %w", query.Error)
	}

	return count > 0, nil
}

func (r *Repository) TrackChannel(channelId string) error {
	query := r.db.
		Clauses(clause.OnConflict{
			OnConstraint: "tracked_channels_pk",
			DoNothing:    true,
		}).
		Table("tracked_channels").
		Create(&TrackedChannel{
			ChannelId: channelId,
		})

	return query.Error
}

func (r *Repository) SaveAttempt(attempt Attempt) error {
	return r.db.
		Clauses(clause.OnConflict{
			OnConstraint: "attempts_pk",
			DoNothing:    true,
		}).
		Table("attempts").
		Create(&attempt).Error
}

func NewRepository() (*Repository, error) {
	var err error
	var repo Repository

	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
		os.Getenv("WORDLE_DB_HOST"),
		os.Getenv("WORDLE_DB_USER"),
		os.Getenv("WORDLE_DB_PASSWORD"),
		os.Getenv("WORDLE_DB_NAME"),
		os.Getenv("WORDLE_DB_PORT"))

	repo.db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("getting database session: %w", err)
	}

	return &repo, nil
}
