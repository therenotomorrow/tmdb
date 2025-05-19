package config

import (
	"context"
	"time"

	"github.com/joho/godotenv"
	"github.com/samber/oops"
	"github.com/sethvargo/go-envconfig"

	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

const (
	service = "config.Settings"

	host    = "https://api.themoviedb.org"
	timeout = 10 * time.Second

	errInvalidConfig = "invalidConfig"
)

type Settings struct {
	Debug bool   `env:"TMDB_DEBUG" json:"debug"`
	Token string `env:"TMDB_TOKEN" json:"token"`

	tmdb.Config
}

func New(filenames ...string) (Settings, error) {
	var settings Settings

	errBuilder := oops.In(service).Code(errInvalidConfig).Public("Invalid configuration.")

	_ = godotenv.Load(filenames...)

	err := envconfig.Process(context.Background(), &settings)
	if err != nil {
		return settings, errBuilder.Wrap(err)
	}

	settings.Config = tmdb.Config{
		Debug:   settings.Debug,
		Timeout: timeout,
		Host:    host,
		Token:   settings.Token,
	}

	return settings, nil
}

func (s *Settings) SetToken(token string) {
	s.Token = token
	s.Config.Token = token
}
