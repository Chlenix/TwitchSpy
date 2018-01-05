package main

import (
	"TwitchSpy/tserror"
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"fmt"
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
	db.Connect("postgres", true)
	defer db.Close()

	tClient := twitch.New()
	defer tClient.RevokeToken()

	topGames, err := tClient.GetTopGames(tClient.Config.Games)
	if err != nil {
		errorTable.handle(*err.(*tserror.AppError))
	}

	fmt.Println(topGames)

	for _, game := range topGames {
		topStreams, _ := tClient.GetStreams(game.GameID)
		fmt.Println(topStreams)
	}
}
