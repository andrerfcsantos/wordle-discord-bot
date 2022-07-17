package db

type UserScore struct {
	UserName string  `gorm:"column:user_name"`
	Score    float64 `gorm:"column:score"`
}

type LeaderboardEntry struct {
	Username    string  `gorm:"column:user_name"`
	TotalScore  float64 `gorm:"column:total_score"`
	AvgAttempts float64 `gorm:"column:avg_attempts"`
	Played      int     `gorm:"column:played"`
}

func (r *Repository) Leaderboard(channelId string) ([]LeaderboardEntry, error) {
	var l []LeaderboardEntry
	query := r.Database().
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
