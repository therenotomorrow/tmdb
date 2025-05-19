package app

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"

	"github.com/samber/oops"

	"github.com/therenotomorrow/tmdb/internal/config"
	"github.com/therenotomorrow/tmdb/pkg/fp"
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
		settings config.Settings
		client   tmdb.Client
		input    io.Reader
		output   io.Writer
		oops     oops.OopsErrorBuilder
	}

	fetchType string

	fetchFunc func(ctx context.Context, page int) ([]tmdb.Movie, error)
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

func (a *TMDB) Fetch(ctx context.Context, page int, kind string) {
	skip := 0
	scan := bufio.NewScanner(a.input)
	fetcher := a.selectFetchFunc(kind)

	for ; ; scan.Scan() {
		switch scan.Text() {
		case "p", "prev", "-":
			page--
		case "n", "next", "+":
			page++
		case "q", "quit", ".":
			return
		default:
			skip++
		}

		if skip > 1 || a.fetch(ctx, page, fetcher) != nil {
			break
		}
	}
}

func (a *TMDB) fetch(ctx context.Context, page int, fetcher fetchFunc) error {
	var err error

	defer func() {
		if err == nil {
			return
		}

		if a.settings.Debug {
			fp.Silent(fmt.Fprintf(a.output, "%+v\n", err))
		} else {
			fp.Silent(fmt.Fprintf(a.output, "%s\n", err))
		}
	}()

	switch {
	case fetcher == nil:
		err = a.oops.Code(errNotFound).
			Public(`Unknown "-type" value for fetch. Allowed [playing,popular,top,upcoming]`).
			New("invalid type")
	case page < 1:
		err = a.oops.Code(errUnexpected).
			Public("Invalid page: Pages start at 1 and max at 500. They are expected to be an integer.").
			New("page less then zero")
	}

	if err != nil {
		return err
	}

	movies, err := oops.Wrap2(fetcher(ctx, page))
	if err != nil {
		return err
	}

	fp.Silent(fmt.Fprintln(a.output))

	for _, movie := range movies {
		fp.Silent(fmt.Fprintf(
			a.output,
			template,
			movie.Title,
			movie.ReleaseDate,
			movie.VoteCount,
			movie.Popularity,
			movie.Overview,
		))
		fp.Silent(a.input.Read(make([]byte, 1)))
	}

	fp.Silent(fmt.Fprintf(a.output, "Current page is %d, prev/next/quit? ", page))

	return nil
}

func (a *TMDB) selectFetchFunc(kind string) fetchFunc {
	var fetch fetchFunc

	switch fetchType(kind) {
	case fetchTypePlaying:
		fetch = a.client.GetNowPlayingMovies
	case fetchTypePopular:
		fetch = a.client.GetPopularMovies
	case fetchTypeTop:
		fetch = a.client.GetTopRatedMovies
	case fetchTypeUpcoming:
		fetch = a.client.GetUpcomingMovies
	}

	return fetch
}
