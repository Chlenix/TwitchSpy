package main

import (
	"TwitchSpy/tserror"
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"fmt"
	"TwitchSpy/util"
	"github.com/labstack/gommon/log"
)

const (
	MaxWarnings = 5
)

type ErrorTable map[error]int

func (errorTable ErrorTable) handle(e tserror.AppError) {
	switch e.Level {
	case tserror.Warning:
		errorTable[e.ErrorObject]++
		if errorTable[e.ErrorObject] >= MaxWarnings {
			fmt.Errorf("%s\n", e.Error())
		}
		break
	case tserror.Critical:
		panic(e.Error())
	}
}

func main() {
	// not concurrency safe
	errorTable := make(ErrorTable)

	// setup db connection
	db.Connect(true)
	defer db.Close()

	// Create twitch client
	tClient := twitch.New()
	defer tClient.RevokeToken()

	// TODO: Skip and develop Stream Insertion
	if !tClient.Config.Debug {

		topGames, err := tClient.GetTopGames()
		if err != nil {
			errorTable.handle(*err.(*tserror.AppError))
		}

		for _, game := range topGames {
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
}
