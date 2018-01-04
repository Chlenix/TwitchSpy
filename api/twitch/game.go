package twitch

import (
	"TwitchSpy/db"
	"reflect"
)

type GameInfo struct {
	Name        string `json:"name"`
	Popularity  int    `json:"popularity"`
	Gameid      int    `json:"_id"`
	Giantbombid int    `json:"giantbomb_id"`
	Genres      []string
	Aliases     []string
	Brief       string
}

func (g GameInfo) Convert() *db.TwitchGame {
	// In progress!
	v := reflect.ValueOf(g)

	values := make([]interface{}, v.NumField())

	for i := 0; i < v.NumField(); i++ {
		values[i] = v.Field(i).Interface()
	}

	return nil
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
