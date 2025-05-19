package app

import (
	"context"
	"io"
	"os"

	"github.com/samber/oops"

	"github.com/therenotomorrow/tmdb/internal/config"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

const (
	service = "app.TMDB"

	template = `---- %q ----
 * Released: %s
 * Votes: %d
 * Popularity: %.2f
 > %s
`
	fetchTypePlaying  fetchType = "playing"
	fetchTypePopular  fetchType = "popular"
	fetchTypeTop      fetchType = "top"
	fetchTypeUpcoming fetchType = "upcoming"

	errUnexpected = "unexpectedError"
	errNotFound   = "notFound"
)

type (
	TMDB struct {
		oops     oops.OopsErrorBuilder
		client   tmdb.Client
		input    io.Reader
		output   io.Writer
		settings config.Settings
	}

	fetchType string

	FetchFunc func(ctx context.Context, page int) ([]tmdb.Movie, error)
)

func New(settings config.Settings) (*TMDB, error) {
	errBuilder := oops.In(service)

	client, err := tmdb.New(settings.Config)
	if err != nil {
		return nil, errBuilder.Wrap(err)
	}

	return &TMDB{settings: settings, client: client, input: os.Stdin, output: os.Stdout, oops: errBuilder}, nil
}

func (a *TMDB) WithDependencies(deps ...any) *TMDB {
	for _, dep := range deps {
		switch dep := dep.(type) {
		case tmdb.Client:
			a.client = dep
		case io.Reader:
			a.input = dep
		case io.Writer:
			a.output = dep
		}
	}

	return a
}

func (a *TMDB) Close() error {
	return a.oops.Code(errUnexpected).Public("Cannot close application.").Wrap(a.client.Close())
}
