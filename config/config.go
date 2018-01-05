package config

type Config struct {
	Debug   bool `envconfig:"DEBUG"`
	Streams int  `envconfig:"STREAMS_PER_GAME"`
	Games   int  `envconfig:"GAMES_LIMIT"`
}
