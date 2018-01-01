package twitch

import (
	"fmt"
	"TwitchSpy/twitch"
	"github.com/levigross/grequests"
	"TwitchSpy/tserror"
)

//	ClientID = "kd1unb4b3q4t58fwlpcbzcbnm76a8fp"
//	AuthHash = "OAuth ec6qnwqbbsjr0k9wh12xfcyr645opo"

const (
	BaseAPI = "https://api.twitch.tv"

	TopGamesEP = "kraken/games/top"
	OAuthEP    = "kraken/oauth2/token"
	RevokeEP   = "kraken/oauth2/revoke"
	RefreshEP  = "kraken/oauth2/refresh"

	ClientID     = "9f859dj8d5z4ydejdwtl2xjhodopvk"
	ClientSecret = "sx4szdadyu1px8t85ha9mf7kwrwxf6"
)

type AuthResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func New() *Client {
	return &Client{
		Session: grequests.NewSession(&grequests.RequestOptions{
			UserAgent: "TwitchSpy/v0.1",
		}),
		ID:     ClientID,
		Secret: ClientSecret,
	}
}

type Client struct {
	Session      *grequests.Session
	ID           string
	Secret       string
	accessToken  string
	refreshToken string
	expiresIn    int
}

func (client Client) RevokeToken() error {
	revokeUrl := fmt.Sprintf("%s/%s", BaseAPI, RevokeEP)
	resp, err := grequests.Post(revokeUrl, &grequests.RequestOptions{
		Params: map[string]string{
			"client_id": client.ID,
			"token":     client.accessToken,
		},
	})

	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	resp.Close()

	return nil
}

func (client Client) RefreshToken() error {
	refreshUrl := fmt.Sprintf("%s/%s", BaseAPI, RefreshEP)
	fmt.Println(refreshUrl)
	return nil
}

func (client *Client) Authorize() error {
	authUrl := fmt.Sprintf("%s/%s", BaseAPI, OAuthEP)
	resp, err := client.Session.Post(authUrl, &grequests.RequestOptions{
		Params: map[string]string{
			"client_id":     client.ID,
			"client_secret": client.Secret,
			"grant_type":    "client_credentials",
		},
	})

	if err != nil {
		return tserror.New(err, tserror.Critical)
	}

	defer resp.Close()

	var jsonSturct AuthResponse
	if err := resp.JSON(&jsonSturct); err != nil {
		return tserror.New(err, tserror.Critical)
	}

	client.accessToken = jsonSturct.AccessToken
	client.refreshToken = jsonSturct.RefreshToken
	client.expiresIn = jsonSturct.ExpiresIn

	return nil
}

func (client Client) GetTopGames(limit, offset int) ([]twitch.Game, error) {
	//gamesUrl := fmt.Sprintf("%s/%s", BaseAPI, TopGamesEP)
	return nil, nil
}
