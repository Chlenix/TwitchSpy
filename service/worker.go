package service

import (
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"time"
	"github.com/labstack/gommon/log"
)

type Worker struct {
	id         int
	WorkerPool chan chan int // Job not int
	JobChannel chan int      // Job not int
	quit       chan bool
	logger     *log.Logger
	started    bool
}

func PulseTopGames() {

	// setup db connection
	db.Connect(true)
	defer db.Close()

	// Create twitch client
	tClient := twitch.New()
	defer tClient.RevokeToken()

	topGames, err := tClient.GetTopGames()
	if err != nil {
		panic(err)
	}

	theTime := time.Now()

	for i, game := range topGames {

		if !db.GameExists(game.GameID) {

			rowsAffected := db.InsertGame(game.Name, game.GameID, game.GiantBombID)

			if tClient.Config.Debug {
				if rowsAffected == 0 {
					log.Printf("Game %s already exists. Skipping", game.Name)
				} else {
					log.Printf("New game %s added", game.Name)
				}
			}

		}

		// record views
		db.RecordGameStats(game.GameID, i+1, game.Viewers, theTime)
	}
}
