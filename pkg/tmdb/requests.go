package tmdb

import (
	"context"
	"net/http"
	"strconv"

	"resty.dev/v3"
)

type moviesResponse struct {
	Results []Movie `json:"results"`
}

func (c *TMDB) GetNowPlayingMovies(ctx context.Context, page int) ([]Movie, error) {
	return c.get(ctx, "/3/movie/now_playing", page)
}

func (c *TMDB) GetPopularMovies(ctx context.Context, page int) ([]Movie, error) {
	return c.get(ctx, "/3/movie/popular", page)
}

func (c *TMDB) GetTopRatedMovies(ctx context.Context, page int) ([]Movie, error) {
	return c.get(ctx, "/3/movie/top_rated", page)
}

func (c *TMDB) GetUpcomingMovies(ctx context.Context, page int) ([]Movie, error) {
	return c.get(ctx, "/3/movie/upcoming", page)
}

func (c *TMDB) get(ctx context.Context, path string, page int) ([]Movie, error) {
	var data moviesResponse

	resp, err := c.engine.R().
		SetContext(ctx).
		SetResult(&data).
		SetQueryParam("page", strconv.Itoa(page)).
		Get(path)

	return data.Results, c.parseResponse(resp, err)
}

func (c *TMDB) parseResponse(resp *resty.Response, err error) error {
	errBuilder := c.oops.Code(errResponse)

	if err != nil {
		return errBuilder.Public("Cannot fetch data from API.").Wrap(err)
	}

	if status := resp.StatusCode(); status != http.StatusOK {
		err, _ := resp.Error().(*errorResponse)

		return errBuilder.With("status", status).Public(err.StatusMessage).New("invalid response")
	}

	return nil
}
