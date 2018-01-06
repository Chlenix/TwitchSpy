package twitch

import "time"

type Stream struct {
	ID        string    `json:"id"`
	UserID    string    `json:"user_id"`
	StartedAt time.Time `json:"started_at"`

	TimeAccessed time.Time
	Title        string `json:"title"`
	GameID       string `json:"game_id"`
	ViewerCount  int    `json:"viewer_count"`
}

type Streams struct {
	Top []Stream `json:"data"`
}
