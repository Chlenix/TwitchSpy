package main

import (
	"TwitchSpy/tserror"
	"TwitchSpy/db"
	"TwitchSpy/api/twitch"
	"fmt"
	"database/sql"
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

	db.InsertGame(&db.TwitchGame{
		Name: "League of Lols",
		Gameid: 251225,
		Giantbombid: sql.NullInt64{2515, true},
		Brief: sql.NullString{},
		Genres: []sql.NullString{},
		Aliases: []sql.NullString{},
	})

	twitchCrawler := twitch.New()
	defer twitchCrawler.RevokeToken()

	topGames, err := twitchCrawler.GetTopGames(20)
	if err != nil {
		errorTable.handle(*err.(*tserror.AppError))
	}

	fmt.Println(topGames)

	// Remove if false
	//if 1 == 0 {
	//	gbClient := giantbomb.GBClient{}
	//
	//	for i := range topGames {
	//		if topGames[i].GameInfo.Giantbombid != 0 {
	//			if err := gbClient.GetGameInfo(&topGames[i], nil); err != nil {
	//				fmt.Fprintf(os.Stderr, "Error! Message: %s | Level: %d\n", err.Error(), err.Level)
	//			} else {
	//				fmt.Println(topGames[i].Brief)
	//			}
	//		}
	//	}
	//}
}
