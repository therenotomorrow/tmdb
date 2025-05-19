package main

import (
	"context"
	"flag"

	"github.com/therenotomorrow/tmdb/internal/app"
	"github.com/therenotomorrow/tmdb/internal/config"
	"github.com/therenotomorrow/tmdb/pkg/fp"
)

var TMDBToken string //nolint:gochecknoglobals // for opportunity to set via `ldflags`

func args() (int, string) {
	page := flag.Int("page", 1, "Page number")
	kind := flag.String("type", "", "The type of list [playing,popular,top,upcoming]")

	flag.Parse()

	return *page, *kind
}

func main() {
	ctx := context.Background()
	page, kind := args()

	settings := fp.Must(config.New())
	if settings.Token == "" {
		settings.SetToken(TMDBToken)
	}

	tmdb := fp.Must(app.New(settings))

	defer func() { _ = tmdb.Close() }()

	tmdb.Fetch(ctx, page, kind)
}
