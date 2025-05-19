package app

import (
	"bufio"
	"context"
	"fmt"

	"github.com/samber/oops"

	"github.com/therenotomorrow/tmdb/pkg/fp"
)

func (a *TMDB) Select(kind string) FetchFunc {
	var fetch FetchFunc

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

func (a *TMDB) Fetch(ctx context.Context, page int, kind string) {
	skip := 0
	scan := bufio.NewScanner(a.input)
	fetcher := a.Select(kind)

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

func (a *TMDB) fetch(ctx context.Context, page int, fetcher FetchFunc) error {
	var err error

	defer func() {
		if err == nil {
			return
		}

		if a.settings.Debug {
			fp.Silent(fmt.Fprintf(a.output, "%+v\n", err))
		} else {
			fp.Silent(fmt.Fprintf(a.output, "%s\n", oops.GetPublic(err, "Something went wrong.")))
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
