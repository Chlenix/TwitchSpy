package twitch

type GameInfo struct {
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	GameID      int    `json:"_id"`
	GiantBombID int    `json:"giantbomb_id"`
	Genres      []string
	Aliases     []string
	Brief       string
}

type Game struct {
	GameInfo     `json:"game"`
	Viewers  int `json:"viewers"`
	Channels int `json:"channels"`
}

type TopGames struct {
	Total int    `json:"_total"`
	Top   []Game `json:"top"`
}
