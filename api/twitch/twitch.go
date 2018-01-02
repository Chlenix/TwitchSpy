package twitch

// https://tmi.twitch.tv/group/user/<TWITCHUSER>/chatters
// Can take up a long time to answer

import (
	"fmt"
	"github.com/levigross/grequests"
	"TwitchSpy/tserror"
	"strconv"
	"errors"
)

const (
	baseURL = "https://api.twitch.tv"

	topGamesEP = "kraken/games/top"
	oAuthEP    = "kraken/oauth2"

	clientID     = "9f859dj8d5z4ydejdwtl2xjhodopvk"
	clientSecret = "sx4szdadyu1px8t85ha9mf7kwrwxf6"

	v5Accept = "application/vnd.twitchtv.v5+json"
)

type errorResponse struct {
	Error string
	Status int
	Message string
}

type AuthHeaders struct {
	ClientID      string
	Authorization string
}

type Token struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	expired      bool
}

type Client struct {
	token   Token
	headers AuthHeaders
}

func New() *Client {
	return &Client{}
}

func (client *Client) Auth() error {
	oAuthUrl := fmt.Sprintf("%s/%s", baseURL, oAuthEP)
	params := map[string]string{
		"client_id":     clientID,
		"client_secret": clientSecret,
	}

	// if no access token in config/settings or access token expired

	params["grant_type"] = "client_credentials"
	resp, err := grequests.Post(fmt.Sprintf("%s/token", oAuthUrl), &grequests.RequestOptions{
		Params: params,
	})

	if err != nil {
		return tserror.New(err, tserror.Critical)
	}

	defer resp.Close()

	if err := resp.JSON(&client.token); err != nil {
		return tserror.New(err, tserror.Critical)
	}

	client.headers = AuthHeaders{
		Authorization: fmt.Sprintf("Bearer %s", client.token.AccessToken),
		ClientID:      clientID,
	}

	return nil
}

func (client *Client) RevokeToken() error {
	revokeUrl := fmt.Sprintf("%s/%s/revoke", baseURL, oAuthEP)
	resp, err := grequests.Post(revokeUrl, &grequests.RequestOptions{
		Params: map[string]string{
			"client_id": client.headers.ClientID,
			"token":     client.token.AccessToken,
		},
	})

	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	defer resp.Close()

	client.token.expired = true

	return nil
}

func (client Client) GetTopGames(limit int) ([]Game, error) {
	// Top Games URL
	gamesUrl := fmt.Sprintf("%s/%s", baseURL, topGamesEP)

	// Configure Request Options
	ro := grequests.RequestOptions{
		Headers: map[string]string{
			"Accept":    v5Accept,
			"Client-ID": client.headers.ClientID,
		},
		Params: map[string]string {
			"limit": strconv.Itoa(limit),
		},
	}

	resp, err := grequests.Get(gamesUrl, &ro)

	if err != nil {
		return nil, tserror.New(err, tserror.Critical)
	}

	defer resp.Close()

	if !resp.Ok {
		var e errorResponse
		if err := resp.JSON(&e); err != nil {
			return nil, tserror.New(err, tserror.Critical)
		}

		errorMsg := fmt.Sprintf("%v: %v. %v", e.Status, e.Error, e.Message)
		return nil, tserror.New(errors.New(errorMsg), tserror.Critical)
	}

	for k, v := range resp.Header {
		fmt.Printf("\"%v\": \"%v\"\n", k, v[0])
	}

	var games TopGames
	if err := resp.JSON(&games); err != nil {
		return nil, tserror.New(err, tserror.Critical)
	}

	return games.Top, nil
}

func isTokenExpired(token string) bool {
	return false
}
