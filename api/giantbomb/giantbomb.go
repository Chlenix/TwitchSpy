package giantbomb

import (
	"fmt"
	"github.com/levigross/grequests"
	"bytes"
	"encoding/csv"
	"strings"
	"io"
	"log"
	"TwitchSpy/tserror"
	"errors"
	"TwitchSpy/api/twitch"
)

const (
	APIKey     = "6d7160493970c1b963963c7791af32352ac31d91"
	BaseAPIUrl = "http://server.giantbomb.com/api"
)

type GBClient struct {
	Suspended bool
}

func formatFilters(filters []string) string {
	var buffer bytes.Buffer
	var sep = false

	for _, filter := range filters {
		if sep {
			buffer.WriteString(",")
		}
		buffer.WriteString(filter)
		sep = true
	}

	return buffer.String()
}

func (c *GBClient) GetGameInfo(game *twitch.Game, filters []string) *tserror.AppError {

	if c.Suspended {
		return tserror.New(errors.New("GBClient suspended"), tserror.Ignore)
	}

	gameInfoUrl := fmt.Sprintf("%s/game/%d", BaseAPIUrl, game.GiantBombID)

	if filters == nil {
		filters = []string{"aliases", "deck", "genres"}
	}

	params := map[string]string{
		"format":     "json",
		"api_key":    APIKey,
		"field_list": formatFilters(filters),
	}

	ro := grequests.RequestOptions{Params: params}

	resp, err := grequests.Get(gameInfoUrl, &ro)
	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	defer resp.Close()

	if !resp.Ok {
		errmsg := fmt.Sprintf("Giantbomb #%d status code: %d",
			game.GiantBombID, resp.StatusCode)
		return tserror.New(errors.New(errmsg), tserror.Warning)
	}

	var gbResponse GBGameInfo
	err = resp.JSON(&gbResponse)

	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	switch gbResponse.StatusCode {
	case Ok:
		break
	case ObjectNotFound:
	case FilterError:
	case SubscriberOnlyVideo:
	case BadUrlFormat:
		return tserror.New(errors.New(gbResponse.Error), tserror.Warning)
	case InvalidAPIKey:
	default:
		return tserror.New(errors.New(gbResponse.Error), tserror.Suspend)
	}

	game.Brief = gbResponse.Deck

	game.Aliases = strings.Split(gbResponse.Aliases, "\n")
	for _, value := range gbResponse.Genres {
		game.Genres = append(game.Genres, value.Name)
	}

	return nil
}
