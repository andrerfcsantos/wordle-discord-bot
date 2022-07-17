package db

import (
	"gorm.io/gorm/clause"
	"time"
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
	HardMode     bool      `gorm:"column:hard_mode"`
}

func (r *Repository) SaveAttempt(attempt Attempt) error {
	return r.Database().
		Clauses(clause.OnConflict{
			UpdateAll: true,
		}).
		Table("attempts").
		Create(&attempt).Error
}

func (r *Repository) AttemptForMessage(channelId string, messageId string) (*Attempt, error) {
	a := Attempt{}
	query := r.Database().
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
	query := r.Database().
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
