package app_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/therenotomorrow/tmdb/internal/app/mocks"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

func movies() []tmdb.Movie {
	return []tmdb.Movie{{
		Title:       "title1",
		Overview:    "overview1",
		ReleaseDate: "releaseDate1",
		Popularity:  12.34,
		VoteCount:   100,
	}, {
		Title:       "title2",
		Overview:    "overview2",
		ReleaseDate: "releaseDate2",
		Popularity:  56.78,
		VoteCount:   200,
	}}
}

func TestTMDBSelect(t *testing.T) {
	t.Parallel()

	type args struct {
		kind string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "playing", args: args{kind: "playing"}},
		{name: "popular", args: args{kind: "popular"}},
		{name: "top", args: args{kind: "top"}},
		{name: "upcoming", args: args{kind: "upcoming"}},
		{name: "invalid", args: args{kind: "invalid"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			_ = New().Select(tt.args.kind)
		})
	}
}

func TestTMDBFetchInvalidKind(t *testing.T) {
	t.Parallel()

	output := new(strings.Builder)

	obj := New().WithDependencies(output)

	obj.Fetch(t.Context(), 1, "invalid")

	assert.Equal(t, "Unknown \"-type\" value for fetch. Allowed [playing,popular,top,upcoming]\n", output.String())
}

func TestTMDBFetchInvalidPage(t *testing.T) {
	t.Parallel()

	output := new(strings.Builder)

	obj := New().WithDependencies(output)

	obj.Fetch(t.Context(), -1, "top")

	assert.Equal(
		t,
		"Invalid page: Pages start at 1 and max at 500. They are expected to be an integer.\n",
		output.String(),
	)
}

func TestTMDBFetchSuccessAllowedKinds(t *testing.T) {
	t.Parallel()

	want := `
---- "title1" ----
 * Released: releaseDate1
 * Votes: 100
 * Popularity: 12.34
 > overview1
---- "title2" ----
 * Released: releaseDate2
 * Votes: 200
 * Popularity: 56.78
 > overview2
Current page is 1, prev/next/quit? `

	type args struct {
		kind   string
		method string
	}

	tests := []struct {
		name string
		args args
	}{
		{name: "playing", args: args{kind: "playing", method: "GetNowPlayingMovies"}},
		{name: "popular", args: args{kind: "popular", method: "GetPopularMovies"}},
		{name: "top", args: args{kind: "top", method: "GetTopRatedMovies"}},
		{name: "upcoming", args: args{kind: "upcoming", method: "GetUpcomingMovies"}},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output := new(strings.Builder)
			client := mocks.NewMockClient(t)
			client.On(test.args.method, mock.Anything, 1).Return(movies(), nil)

			obj := New().WithDependencies(output, client)

			obj.Fetch(t.Context(), 1, test.args.kind)

			assert.Equal(t, want, output.String())
		})
	}
}
