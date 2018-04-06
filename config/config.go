package config

import "fmt"

type GiantbombConfig struct {
	GBAPIKey string `envconfig:"APIKEY"`
}

type ServerConfig struct {
	Secure   bool   `envconfig:"SECURE"`
	HostName string `envconfig:"HOSTNAME"`
	Port     string `envconfig:"PORT"`
}

type ClientConfig struct {
	Debug            bool `envconfig:"DEBUG"`
	Streams          int  `envconfig:"STREAMS_PER_GAME"`
	GamesToFetch     int  `envconfig:"GAMES_TO_FETCH"`
	WithinTopStreams int  `envconfig:"WITHIN_TOP_STREAMS"`
	MinTopSeconds    int  `envconfig:"MINIMUM_TOP_SECONDS"`
}

type DBConfig struct {
	DBHost string `envconfig:"HOST"`
	DBPort int    `envconfig:"PORT"`
	DBName string `envconfig:"NAME"`
	DBUser string `envconfig:"USER"`
	DBPass string `envconfig:"PASS"`
}

func (dbc DBConfig) ToString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		dbc.DBHost, dbc.DBPort, dbc.DBUser, dbc.DBPass, dbc.DBName, "disable")
}
