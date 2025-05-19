package tmdb

import (
	"context"
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
