package tmdb_test

import (
	"errors"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/tmdb/pkg/tmdb"
	"github.com/therenotomorrow/tmdb/pkg/tmdb/mocks"
)

var errFail = errors.New("fail")

func want() []tmdb.Movie {
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

func TestTMDBGetNowPlayingMoviesSuccess(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(successResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetNowPlayingMovies(t.Context(), 1)

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestTMDBGetNowPlayingMoviesFailure(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(nil, errFail)

	obj := New().SetTransport(trans)
	got, err := obj.GetNowPlayingMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, `Get "https://tmdb.host/3/movie/now_playing?language=en&page=1": fail`)

	assert.Equal(t, "Cannot fetch data from API.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetNowPlayingMoviesError(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(failureResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetNowPlayingMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, "invalid response")

	assert.Equal(t, "Some public message.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetPopularMoviesSuccess(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(successResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetPopularMovies(t.Context(), 1)

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestTMDBGetPopularMoviesFailure(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(nil, errFail)

	obj := New().SetTransport(trans)
	got, err := obj.GetPopularMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, `Get "https://tmdb.host/3/movie/popular?language=en&page=1": fail`)

	assert.Equal(t, "Cannot fetch data from API.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetPopularMoviesError(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(failureResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetPopularMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, "invalid response")

	assert.Equal(t, "Some public message.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetTopRatedMoviesSuccess(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(successResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetTopRatedMovies(t.Context(), 1)

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestTMDBGetTopRatedMoviesFailure(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(nil, errFail)

	obj := New().SetTransport(trans)
	got, err := obj.GetTopRatedMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, `Get "https://tmdb.host/3/movie/top_rated?language=en&page=1": fail`)

	assert.Equal(t, "Cannot fetch data from API.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetTopRatedMoviesError(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(failureResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetTopRatedMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, "invalid response")

	assert.Equal(t, "Some public message.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetUpcomingMoviesSuccess(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(successResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetUpcomingMovies(t.Context(), 1)

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestTMDBGetUpcomingMoviesFailure(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(nil, errFail)

	obj := New().SetTransport(trans)
	got, err := obj.GetUpcomingMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, `Get "https://tmdb.host/3/movie/upcoming?language=en&page=1": fail`)

	assert.Equal(t, "Cannot fetch data from API.", orr.Public())
	assert.Empty(t, got)
}

func TestTMDBGetUpcomingMoviesError(t *testing.T) {
	t.Parallel()

	trans := mocks.NewMockRoundTripper(t)
	trans.On("RoundTrip", mock.Anything).Return(failureResponse(t), nil)

	obj := New().SetTransport(trans)
	got, err := obj.GetUpcomingMovies(t.Context(), 1)

	var orr oops.OopsError

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, "invalid response")

	assert.Equal(t, "Some public message.", orr.Public())
	assert.Empty(t, got)
}

func response(t *testing.T, code int, json string) *http.Response {
	t.Helper()

	resp := new(http.Response)

	resp.StatusCode = code
	resp.Header = make(http.Header)
	resp.Header.Set("Content-Type", "application/json")
	resp.Body = io.NopCloser(strings.NewReader(json))

	return resp
}

func successResponse(t *testing.T) *http.Response {
	t.Helper()

	return response(t, http.StatusOK, `
{
  "results": [
    {
      "title": "title1",
      "overview": "overview1",
      "release_date": "releaseDate1",
      "popularity": 12.34,
      "vote_count": 100
    },
    {
      "title": "title2",
      "overview": "overview2",
      "release_date": "releaseDate2",
      "popularity": 56.78,
      "vote_count": 200
    }
  ]
}
`)
}

func failureResponse(t *testing.T) *http.Response {
	t.Helper()

	return response(t, http.StatusTeapot, `
{
  "status_message": "Some public message.",
  "status_code": 42
}
`)
}
