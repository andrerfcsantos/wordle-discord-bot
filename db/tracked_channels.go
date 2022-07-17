package db

import (
	"fmt"
	"gorm.io/gorm/clause"
)

type TrackedChannel struct {
	ChannelId string `gorm:"primary_key;column:channel_id"`
}

func (r *Repository) IsTrackedChannel(channelId string) (bool, error) {
	var count int64

	query := r.Database().
		Table("tracked_channels").
		Where("channel_id = ?", channelId).
		Count(&count)

	if query.Error != nil {
		return false, fmt.Errorf("checking tracked channel: %w", query.Error)
	}

	return count > 0, nil
}

func (r *Repository) TrackChannel(channelId string) error {
	query := r.Database().
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
