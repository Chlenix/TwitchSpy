package giantbomb

import (
	"TwitchSpy/twitch"
	"fmt"
	"github.com/levigross/grequests"
	"bytes"
	"encoding/csv"
	"strings"
	"io"
	"log"
	"TwitchSpy/tserror"
	"errors"
)

const (
	APIKey  = "6d7160493970c1b963963c7791af32352ac31d91"
	BaseAPI = "http://www.giantbomb.com/api"
)

type GBClient struct {
	RequestOptions *grequests.RequestOptions
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

func parseAliases(in string) []string {
	r := csv.NewReader(strings.NewReader(in))
	r.Comma = '\n'

	var aliases []string

	for {
		alias, err := r.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			log.Fatal(err)
		}

		aliases = append(aliases, alias[0])
	}

	return aliases
}

func (c *GBClient) GetGameInfo(game *twitch.Game, filters []string) *tserror.AppError {
	gameInfoUrl := fmt.Sprintf("%s/game/%d", BaseAPI, game.Giantbombid)

	if filters == nil {
		filters = []string{"aliases", "deck", "genres"}
	}

	c.RequestOptions.Params = map[string]string{
		"format":     "json",
		"api_key":    APIKey,
		"field_list": formatFilters(filters),
	}

	resp, err := grequests.Get(gameInfoUrl, c.RequestOptions)
	if err != nil {
		return tserror.New(err, tserror.Warning)
	}

	defer resp.Close()

	if !resp.Ok {
		errmsg := fmt.Sprintf("Giantbomb #%d status code: %d",
			game.Giantbombid, resp.StatusCode)
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
	game.Aliases = parseAliases(gbResponse.Aliases)
	for _, value := range gbResponse.Genres {
		game.Genres = append(game.Genres, value.Name)
	}

	return nil
}
