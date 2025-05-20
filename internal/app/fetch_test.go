package app_test

import (
	"errors"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/therenotomorrow/tmdb/internal/app/mocks"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

var errFail = errors.New("fail")

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

func TestTMDBFetchFailure(t *testing.T) {
	t.Parallel()

	type args struct {
		debug bool
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{name: "public error", args: args{debug: false}, want: "Something went wrong."},
		{name: "debug error", args: args{debug: true}, want: "Oops: fail"},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			output := new(strings.Builder)
			client := mocks.NewMockClient(t)
			client.On("GetTopRatedMovies", mock.Anything, 2).Return(nil, errFail)

			obj := New(test.args.debug).WithDependencies(output, client)

			obj.Fetch(t.Context(), 2, "top")

			assert.Contains(t, output.String(), test.want)
		})
	}
}

type reader struct {
	data []string
	curr int
}

func newReader(data ...string) *reader {
	return &reader{data: data, curr: 0}
}

func (r *reader) Read(buf []byte) (int, error) {
	if len(buf) == 1 {
		return 1, nil
	}

	if r.curr >= len(r.data) {
		return 0, io.EOF
	}

	data := r.data[r.curr] + "\n"

	copy(buf, data)

	r.curr++

	return len(data), nil
}

func TestTMDBFetchSuccessFlow(t *testing.T) {
	t.Parallel()

	input := newReader("next", "prev", "quit", "next")
	output := new(strings.Builder)
	client := mocks.NewMockClient(t)
	client.On("GetTopRatedMovies", mock.Anything, 1).Times(2).Return(movies(), nil)
	client.On("GetTopRatedMovies", mock.Anything, 2).Times(1).Return(movies(), nil)

	obj := New().WithDependencies(input, output, client)

	obj.Fetch(t.Context(), 1, "top")
}

func TestTMDBFetchSuccessSkip(t *testing.T) {
	t.Parallel()

	input := newReader("next", "skip", "prev", "quit")
	output := new(strings.Builder)
	client := mocks.NewMockClient(t)
	client.On("GetTopRatedMovies", mock.Anything, 1).Times(1).Return(movies(), nil)
	client.On("GetTopRatedMovies", mock.Anything, 2).Times(1).Return(movies(), nil)

	obj := New().WithDependencies(input, output, client)

	obj.Fetch(t.Context(), 1, "top")
}

func TestTMDBFetchFailureFlow(t *testing.T) {
	t.Parallel()

	input := newReader("next", "skip")
	output := new(strings.Builder)
	client := mocks.NewMockClient(t)
	client.On("GetTopRatedMovies", mock.Anything, 1).Times(1).Return(movies(), nil)
	client.On("GetTopRatedMovies", mock.Anything, 2).Times(1).Return(nil, errFail)

	obj := New().WithDependencies(input, output, client)

	obj.Fetch(t.Context(), 1, "top")
}
