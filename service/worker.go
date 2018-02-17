package service

import (
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"time"
	"TwitchSpy/util"
	"github.com/labstack/gommon/log"
	"fmt"
)

type Worker struct {
	id         int
	WorkerPool chan chan int // Job not int
	JobChannel chan int      // Job not int
	quit       chan bool
	logger     *log.Logger
	started    bool
}

func ViewStream(id int, in <-chan int) {
	defer wg.Done()
	var totalChunks = 0

	fmt.Printf("Viewer [%v]: initialized\n", id)

	for {
		_, more := <-in

		if !more {
			break
		}

		totalChunks++
	}

	fmt.Printf("Viewer [%v]: Downloaded a total of %d video chunks. EXITED\n", id, totalChunks)
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

			g := &db.TwitchGame{
				Name:        game.Name,
				GameID:      game.GameID,
				GiantBombID: util.ToNullInt64(game.GiantBombID),
			}

			rowsAffected := db.InsertGame(g)

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
