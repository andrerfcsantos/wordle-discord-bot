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

type LeaderboardEntry struct {
	Username    string  `gorm:"column:user_name"`
	AvgScore    float64 `gorm:"column:avg_score"`
	AvgAttempts float64 `gorm:"column:avg_attempts"`
	TotalPoints int     `gorm:"column:total_points"`
	Count       int     `gorm:"column:count"`
}

func (r *Repository) Leaderboard(channelId string) ([]LeaderboardEntry, error) {
	l := []LeaderboardEntry{}
	query := r.db.
		Raw(`
		select
			a.user_name, trunc(avg(a.score), 3) as "avg_score", trunc(avg(a.attempts), 3) "avg_attempts", sum(score) "total_points", count(*) as "count"
		from 
			attempts a
		where 
			a.channel_id = ?
		group by
			a.user_name
		order by 
			2 desc;
		`, channelId).
		Scan(&l)

	return l, query.Error
}

type UserScore struct {
	UserName string  `gorm:"column:user_name"`
	Score    float64 `gorm:"column:score"`
}

func (r *Repository) LeaderboardForDay(channelId string, day int) ([]UserScore, error) {
	l := []UserScore{}
	query := r.db.
		Raw(`
		select
			a.user_name, score
		from 
			attempts a
		where 
			a."day" = ? and a.channel_id = ?
		order by 
			2 desc;`, day, channelId).
		Scan(&l)

	return l, query.Error
}

func (r *Repository) SaveAttempt(attempt Attempt) error {
	return r.db.
		Clauses(clause.OnConflict{
			UpdateAll: true,
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
