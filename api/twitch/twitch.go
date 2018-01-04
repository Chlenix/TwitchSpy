package twitch

// https://tmi.twitch.tv/group/user/<TWITCHUSER_ALL_LOWERCASE>/chatters
// Can take up a long time to answer

import (
	"fmt"
	"github.com/levigross/grequests"
	"TwitchSpy/tserror"
	"strconv"
	"errors"
	"TwitchSpy/db"
)

const (
	baseURL = "https://api.twitch.tv"

	topGamesEP = "kraken/games/top"
	oAuthEP    = "kraken/oauth2"

	v5Accept = "application/vnd.twitchtv.v5+json"
)

type errorResponse struct {
	Error   string
	Status  int
	Message string
}

type authHeaders struct {
	ClientSecret  string
	ClientID      string
	Authorization string
}

type ClientToken struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expired      bool
}

type Client struct {
	Token   ClientToken
	Headers authHeaders
}

// Returns a new client object with state information about
// the client token and pre-set headers
func New() *Client {
	var c Client

	var meta = db.GetClient()

	c.Headers.ClientID = meta.ClientID
	c.Headers.ClientSecret = meta.ClientSecret
	c.Token.AccessToken = meta.AccessToken.String
	c.Token.RefreshToken = meta.RefreshToken.String
	c.Token.Expired = meta.Expired

	// Auth request new token if current has expired or empty token
	if c.Token.Expired || !meta.AccessToken.Valid {
		if err := c.Auth(); err != nil {
			panic(err)
		}
	}

	c.Headers.Authorization = fmt.Sprintf("Bearer %s", c.Token.AccessToken)

	return &c
}

func (client *Client) Auth() error {

	// Preparing the oAuth endpoint url and the required GET parameters
	oAuthUrl := fmt.Sprintf("%s/%s", baseURL, oAuthEP)
	params := map[string]string{
		"client_id":     client.Headers.ClientID,
		"client_secret": client.Headers.ClientSecret,
		"grant_type": "client_credentials",
	}

	ro := grequests.RequestOptions{Params: params}

	resp, err := grequests.Post(fmt.Sprintf("%s/token", oAuthUrl), &ro)

	if err != nil {
		return tserror.New(err, tserror.Critical)
	}

	defer resp.Close()

	if err := resp.JSON(&client.Token); err != nil {
		return tserror.New(err, tserror.Critical)
	}

	// Update the database token information
	client.Token.Expired = false
	return db.UpdateClientToken(db.ClientToken(client.Token))
}

func (client *Client) RevokeToken() error {
	revokeUrl := fmt.Sprintf("%s/%s/revoke", baseURL, oAuthEP)
	resp, err := grequests.Post(revokeUrl, &grequests.RequestOptions{
		Params: map[string]string{
			"client_id": client.Headers.ClientID,
			"token":     client.Token.AccessToken,
		},
	})

	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	defer resp.Close()

	client.Token.Expired = true
	db.UpdateClientToken(db.ClientToken(client.Token))

	return nil
}

func (client Client) GetTopGames(limit int) ([]Game, error) {
	// Top Games URL
	gamesUrl := fmt.Sprintf("%s/%s", baseURL, topGamesEP)

	// Configure Request Options
	ro := grequests.RequestOptions{
		Headers: map[string]string{
			"Accept":    v5Accept,
			"Client-ID": client.Headers.ClientID,
		},
		Params: map[string]string{
			"limit": strconv.Itoa(limit),
		},
	}

	resp, err := grequests.Get(gamesUrl, &ro)

	if err != nil {
		return nil, tserror.New(err, tserror.Critical)
	}

	defer resp.Close()

	// If the response is not 2xx (OK)
	if !resp.Ok {
		var e errorResponse
		if err := resp.JSON(&e); err != nil {
			return nil, tserror.New(err, tserror.Critical)
		}

		errorMsg := fmt.Sprintf("%v: %v. %v", e.Status, e.Error, e.Message)
		return nil, tserror.New(errors.New(errorMsg), tserror.Critical)
	}

	var games TopGames
	if err := resp.JSON(&games); err != nil {
		return nil, tserror.New(err, tserror.Critical)
	}

	return games.Top, nil
}
