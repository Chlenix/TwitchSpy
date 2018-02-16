package service

import (
	"TwitchSpy/api/twitch"
	"TwitchSpy/config"
)

type Service struct {
	Dispatcher *Dispatcher
	Client     *twitch.Client
	MainConfig *config.ClientConfig
	DBConfig   *config.DBConfig
}
