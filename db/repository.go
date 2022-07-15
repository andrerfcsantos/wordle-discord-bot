package db

import (
	"fmt"
	log "github.com/sirupsen/logrus"
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
	HardMode     bool      `gorm:"column:hard_mode"`
}

type TrackedChannel struct {
	ChannelId string `gorm:"primary_key;column:channel_id"`
}

type Repository struct {
	db *gorm.DB
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
	TotalScore  float64 `gorm:"column:total_score"`
	AvgAttempts float64 `gorm:"column:avg_attempts"`
	Played      int     `gorm:"column:played"`
}

func (r *Repository) Leaderboard(channelId string) ([]LeaderboardEntry, error) {
	l := []LeaderboardEntry{}
	query := r.db.
		Raw(`
		select
			a.user_name,
			sum(a.new_score) as "total_score",
			avg(a.attempts) as "avg_attempts",
			count(*) as "played"
		from (
			select
				user_name,
				attempts,
				(30 - ((select date_part as value from DATE_PART('day', now() - '2021-06-19')) - day) )/30 * (
				case
					when a.success then (7-attempts)*(7-attempts)+2
					else 2
				end) new_score
			from
				attempts a
			where
				channel_id = ? and 
				day > (DATE_PART('day', now() - '2021-06-19') -30)
			order by day desc
		) a
		group by a.user_name
		order by 2 desc
		`, channelId).
		Scan(&l)

	return l, query.Error
}

type UserScore struct {
	UserName string  `gorm:"column:user_name"`
	Score    float64 `gorm:"column:score"`
}

func (r *Repository) AttemptForMessage(channelId string, messageId string) (*Attempt, error) {
	a := Attempt{}
	query := r.db.
		Raw(`
		select
			*
		from 
			attempts a
		where 
			a.channel_id = ? and a.message_id = ?;`,
			channelId, messageId).Scan(&a)

	return &a, query.Error
}

func (r *Repository) DeleteAttemptForMessage(channelId string, messageId string) (bool, error) {
	query := r.db.
		Exec(`
		delete from
			attempts a
		where 
			a.channel_id = ? and a.message_id = ?;`,
			channelId, messageId)

	var deleted bool
	if query.RowsAffected != 0 {
		deleted = true
	}

	return deleted, query.Error
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

func NewRepository() (*Repository, error) {
	var err error
	var repo Repository

	repo.db, err = getConnection()
	if err != nil {
		return nil, fmt.Errorf("getting database session: %w", err)
	}

	return &repo, nil
}
