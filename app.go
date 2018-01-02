package main

import (
	"TwitchSpy/api/twitch"
	"fmt"
	//"TwitchSpy/api/giantbomb"
	//"github.com/levigross/grequests"
	//"os"
	"TwitchSpy/tserror"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/jmoiron/sqlx"
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
		}
		break
	case tserror.Critical:
		break
	}
}

func main() {

	// not concurrency safe
	errorTable := make(ErrorTable)

	twitchCrawler := twitch.New()
	if err := twitchCrawler.Auth(); err != nil {
		errorTable.handle(*err.(*tserror.AppError))
	}

	defer twitchCrawler.RevokeToken()

	topGames, err := twitchCrawler.GetTopGames(20)
	if err != nil {
		errorTable.handle(*err.(*tserror.AppError))
	}

	fmt.Println(topGames)

	//var gbClient = &giantbomb.GBClient{
	//	RequestOptions: &grequests.RequestOptions{},
	//}
	//
	//for i := range topGames {
	//	if topGames[i].GameInfo.Giantbombid != 0 {
	//		if err := gbClient.GetGameInfo(&topGames[i], nil); err != nil {
	//			fmt.Fprintf(os.Stderr, "Error! Message: %s | Level: %d\n", err.Error(), err.Level)
	//		} else {
	//			fmt.Println(topGames[i].Brief)
	//		}
	//	}
	//}
}
