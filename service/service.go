package service

import (
	"TwitchSpy/api/twitch"
	"TwitchSpy/config"
)

type Service struct {
	Dispatcher *Dispatcher
	Listener   *Listener
	Client     *twitch.Client
	MainConfig *config.ClientConfig
	DBConfig   *config.DBConfig
}

func NewService() *Service {
	return &Service{
		Dispatcher: nil,
		Listener:   nil,
		Client:     twitch.New(),
		MainConfig: nil,
		DBConfig:   nil,
	}
}
