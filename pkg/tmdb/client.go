package tmdb

import (
	"context"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/samber/oops"
	"resty.dev/v3"
)

const (
	service = "tmdb.TMDB"

	outOfRangeCode = 22

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
		Debug   bool
		Timeout time.Duration `validate:"required,min=5s"`
		Host    string        `validate:"required,url"`
		Token   string        `validate:"required"`
	}

	TMDB struct {
		config Config
		engine *resty.Client
		oops   oops.OopsErrorBuilder
	}

	errorResponse struct {
		StatusCode    int    `json:"status_code"`
		StatusMessage string `json:"status_message"`
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

func (c *TMDB) Close() error {
	return c.oops.Code(errUnexpected).Public("Cannot close client connections.").Wrap(c.engine.Close())
}

func (c *TMDB) WithDependencies(rt http.RoundTripper) *TMDB {
	c.engine = c.engine.SetTransport(rt)

	return c
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
	errBuilder := c.oops.Code(errResponse).Public("Cannot get now playing movies.")

	if err != nil {
		return errBuilder.Wrap(err)
	}

	if status := resp.StatusCode(); status != http.StatusOK {
		err, ok := resp.Error().(*errorResponse)
		if !ok {
			err = unexpected()
		}

		if err.StatusCode == outOfRangeCode {
			errBuilder = errBuilder.Public("Page number out of range.")
		}

		return errBuilder.With("status", status).New(err.StatusMessage)
	}

	return nil
}

func unexpected() *errorResponse {
	return &errorResponse{
		StatusCode:    -1,
		StatusMessage: "Unexpected response from the API.",
	}
}
