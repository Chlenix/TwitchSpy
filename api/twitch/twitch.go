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
	"TwitchSpy/config"
	"github.com/kelseyhightower/envconfig"
)

const (
	baseURL = "https://api.twitch.tv"

	topGamesEP = "kraken/games/top"
	oAuthEP    = "kraken/oauth2"

	streamsEP = "helix/streams"
	v5Accept  = "application/vnd.twitchtv.v5+json"
)

type errorResponse struct {
	Error   string
	Status  int
	Message string
}

type StreamQueryOptions struct {
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
	Config  *config.ClientConfig
}

// Returns a new client object with state information about
// the client token and pre-set headers
func New() *Client {
	var c Client
	c.Config = &config.ClientConfig{}
	// Populate the environment variables
	if err := envconfig.Process("ts", c.Config); err != nil {
		panic(err)
	}

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
		"grant_type":    "client_credentials",
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

// Gets information about active streams. Streams are returned sorted by number of current viewers, in descending order.
// Across multiple pages of results, there may be duplicate or missing streams, as viewers join and leave streams.
// The response has a JSON payload with a data field containing an array of stream information elements
// and a pagination field containing information required to query for more streams.
func (client Client) GetStreams(gameID int) ([]Stream, error) {
	streamsUrl := fmt.Sprintf("%s/%s", baseURL, streamsEP)

	var ro grequests.RequestOptions
	ro.Headers = map[string]string{
		"Client-ID":     client.Headers.ClientID,
		"Authorization": client.Headers.Authorization,
	}

	ro.Params = map[string]string{
		// Maximum number of objects to return. Maximum: 100. Default: 20.
		"first": strconv.Itoa(client.Config.Streams),

		// 	Stream language. You can specify up to 100 languages.
		"language": "en",

		// Returns streams broadcasting a specified game ID. You can specify up to 100 IDs.
		"game_id": strconv.Itoa(gameID),

		// Stream type: "all", "live", "vodcast". Default: "all".
		"type": "live",
	}

	resp, err := grequests.Get(streamsUrl, &ro)
	if err != nil {
		panic(err)
	}

	defer resp.Close()

	if !resp.Ok {
		panic(err)
	}

	var topStreams Streams
	if err := resp.JSON(&topStreams); err != nil {
		panic(err)
	}

	fmt.Println(topStreams.Top[0].StartedAt.String())

	return topStreams.Top, nil
}

func (client Client) GetTopGames() ([]Game, error) {
	// Top GamesToFetch URL
	gamesUrl := fmt.Sprintf("%s/%s", baseURL, topGamesEP)

	// Configure Request Options
	ro := grequests.RequestOptions{
		Headers: map[string]string{
			"Accept":    v5Accept,
			"Client-ID": client.Headers.ClientID,
		},
		Params: map[string]string{
			"limit": strconv.Itoa(client.Config.GamesToFetch),
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
