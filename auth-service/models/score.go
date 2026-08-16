package models

import "time"

type Score struct {
	UserID    int
	Score     int
	CreatedAt time.Time
}

type ScoreWithUsername struct {
	Username string `json:"username"`
	Score    int    `json:"score"`
}
