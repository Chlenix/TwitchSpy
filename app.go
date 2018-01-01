package main

import (
	"TwitchSpy/api/twitch"
	"fmt"
	"TwitchSpy/api/giantbomb"
	"github.com/levigross/grequests"
	"os"
)

func main() {

	// 20674 --- 2 genres
	var twitchCrawler = twitch.New()
	twitchCrawler.Authorize()

	twitchCrawler.RevokeToken()

	topGames, err := twitchCrawler.GetTopGames(24, 0)
	if err != nil {
		fmt.Println(err.Error())
	}

	var gbClient = &giantbomb.GBClient{
		RequestOptions: &grequests.RequestOptions{},
	}

	for i := range topGames {
		if topGames[i].GameInfo.Giantbombid != 0 {
			if err := gbClient.GetGameInfo(&topGames[i], nil); err != nil {
				fmt.Fprintf(os.Stderr, "Error! Message: %s | Level: %d\n", err.Error(), err.Level)
			} else {
				fmt.Println(topGames[i].Brief)
			}
		}
	}
}
