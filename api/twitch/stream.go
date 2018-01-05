package twitch

import "time"

type Stream struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	GameID      string    `json:"game_id"`
	ViewerCount int       `json:"viewer_count"`
	StartedAt   time.Time `json:"started_at"`
}

type Streams struct {
	Top []Stream `json:"data"`
}
