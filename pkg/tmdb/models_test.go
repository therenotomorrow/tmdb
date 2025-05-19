package tmdb_test

import (
	"testing"

	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

func TestMovieStructure(t *testing.T) {
	t.Parallel()

	_ = tmdb.Movie{
		Title:       "title",
		Overview:    "overview",
		Popularity:  19.92,
		ReleaseDate: "released",
		VoteCount:   666,
	}
}
