package tmdb

import (
	"context"
	"io"
	"net/http"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"resty.dev/v3"
)

const (
	service = "tmdb.TMDB"

	errInvalidConfig = "invalidConfig"
	errUnexpected    = "unexpectedError"
	errResponse      = "responseError"
)

type (
	Client interface {
		GetNowPlayingMovies(ctx context.Context, page int) ([]Movie, error)
		GetPopularMovies(ctx context.Context, page int) ([]Movie, error)
		GetTopRatedMovies(ctx context.Context, page int) ([]Movie, error)
		GetUpcomingMovies(ctx context.Context, page int) ([]Movie, error)
		io.Closer
	}

	Config struct {
		Host    string        `validate:"required,url"`
		Token   string        `validate:"required"`
		Timeout time.Duration `validate:"required,min=5s"`
		Debug   bool
	}

	TMDB struct {
		oops   oops.OopsErrorBuilder
		engine *resty.Client
		config Config
	}

	errorResponse struct {
		StatusMessage string `json:"status_message"`
		StatusCode    int    `json:"status_code"`
	}
)

func New(config Config) (*TMDB, error) {
	errBuilder := oops.In(service)

	if err := validator.New().Struct(config); err != nil {
		return nil, errBuilder.Code(errInvalidConfig).Public("Invalid configuration.").Wrap(err)
	}

	engine := resty.New().
		SetTimeout(config.Timeout).
		SetAuthToken(config.Token).
		SetBaseURL(config.Host).
		SetDebug(config.Debug).
		SetError(new(errorResponse)).
		SetHeader("accept", "application/json").
		SetQueryParams(map[string]string{"language": "en"})

	return &TMDB{config: config, engine: engine, oops: errBuilder}, nil
}

func (c *TMDB) SetTransport(rt http.RoundTripper) *TMDB {
	c.engine = c.engine.SetTransport(rt)

	return c
}

func (c *TMDB) Close() error {
	return c.oops.Code(errUnexpected).Public("Cannot close client connections.").Wrap(c.engine.Close())
}
