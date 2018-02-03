package main

import (
	"TwitchSpy/tserror"
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"TwitchSpy/util"
	"github.com/labstack/gommon/log"
)

func main() {
	// not concurrency safe
	errorTable := make(tserror.ErrorTable)

	// setup db connection
	db.Connect(true)
	defer db.Close()

	// Create twitch client
	tClient := twitch.New()
	defer tClient.RevokeToken()

	topGames, err := tClient.GetTopGames()
	if err != nil {
		errorTable.Handle(*err.(*tserror.AppError))
	}

	for _, game := range topGames {

		if db.GameExists(game.GameID) {
			log.Printf("Game %s already exists. Skipping", game.Name)
			continue
		}

		// implement cache slice, if present then already set
		rowsAffected := db.InsertGame(&db.TwitchGame{
			Name:        game.Name,
			GameID:      game.GameID,
			GiantBombID: util.ToNullInt64(game.GiantBombID),
		})

		if tClient.Config.Debug {
			if rowsAffected == 0 {
				log.Printf("Game %s already exists. Skipping", game.Name)
			} else {
				log.Printf("New game %s added", game.Name)
			}
		}
	}
}
