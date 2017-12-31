package twitch

import (
	"fmt"
	"TwitchSpy/twitch"
	"github.com/levigross/grequests"
	"TwitchSpy/tserror"
)

// Official Client-ID: 9f859dj8d5z4ydejdwtl2xjhodopvk
// Official Client Secret: sx4szdadyu1px8t85ha9mf7kwrwxf6

//const (
//	BaseAPI = "https://api.twitch.tv"
//	TopGamesEP = "kraken/games/top"
//
//	ClientID = "kd1unb4b3q4t58fwlpcbzcbnm76a8fp"
//	AuthHash = "OAuth ec6qnwqbbsjr0k9wh12xfcyr645opo"
//)

const (
	BaseAPI    = "https://api.twitch.tv"
	TopGamesEP = "kraken/games/top"
	OAuthEP    = "kraken/oauth2/token"

	ClientID     = "9f859dj8d5z4ydejdwtl2xjhodopvk"
	ClientSecret = "sx4szdadyu1px8t85ha9mf7kwrwxf6"
)

func New() *Client {
	return &Client{
		Session: grequests.NewSession(&grequests.RequestOptions{
			UserAgent: "TwitchSpy/v0.1",
		}),
		clientID: ClientID,
		clientSecret: ClientSecret,
	}
}

type Client struct {
	Session      *grequests.Session
	clientID     string
	clientSecret string
	accessToken  string
	refreshToken string
	expiresIn    int
}

func (client *Client) Authorize() error {
	authUrl := fmt.Sprintf("%s/%s", BaseAPI, OAuthEP)
	resp, err := client.Session.Post(authUrl, &grequests.RequestOptions{
		Params: map[string]string {
			"client_id": client.clientID,
			"client_secret": client.clientSecret,
			"grant_type": "client_credentials",
			"scope": "viewing_activity_read",
		},
	})

	if err != nil {
		return tserror.New(err, tserror.Critical)
	}

	ss := resp.String()
	fmt.Println(ss)

	return nil
}

func (client *Client) GetTopGames(limit, offset int) ([]twitch.Game, error) {
	//gamesUrl := fmt.Sprintf("%s/%s", BaseAPI, TopGamesEP)
	return nil, nil
}
