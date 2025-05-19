package tmdb_test

import (
	"net/http"
	"testing"
	"time"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/tmdb/pkg/fp"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
	"github.com/therenotomorrow/tmdb/pkg/tmdb/mocks"
)

func New() *tmdb.TMDB {
	return fp.Must(tmdb.New(tmdb.Config{
		Debug:   false,
		Timeout: time.Minute,
		Host:    "https://tmdb.host",
		Token:   "secret",
	}))
}

func TestClient(t *testing.T) {
	t.Parallel()

	var _ tmdb.Client = (*tmdb.TMDB)(nil)
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		config tmdb.Config
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "min config",
			args: args{config: tmdb.Config{
				Debug:   false,
				Timeout: time.Minute,
				Host:    "https://tmdb.host",
				Token:   "secret",
			}},
			want: "",
		},
		{
			name: "max config",
			args: args{config: tmdb.Config{
				Debug:   true,
				Timeout: time.Minute,
				Host:    "https://tmdb.host",
				Token:   "secret",
			}},
			want: "",
		},
		{
			name: "small timeout",
			args: args{config: tmdb.Config{
				Debug:   true,
				Timeout: time.Second,
				Host:    "https://tmdb.host",
				Token:   "secret",
			}},
			want: "Key: 'Config.Timeout' Error:Field validation for 'Timeout' failed on the 'min' tag",
		},
		{
			name: "missing timeout",
			args: args{config: tmdb.Config{
				Debug:   true,
				Timeout: 0,
				Host:    "https://tmdb.host",
				Token:   "secret",
			}},
			want: "Key: 'Config.Timeout' Error:Field validation for 'Timeout' failed on the 'required' tag",
		},
		{
			name: "missing host",
			args: args{config: tmdb.Config{
				Debug:   true,
				Timeout: time.Minute,
				Host:    "",
				Token:   "secret",
			}},
			want: "Key: 'Config.Host' Error:Field validation for 'Host' failed on the 'required' tag",
		},
		{
			name: "missing token",
			args: args{config: tmdb.Config{
				Debug:   true,
				Timeout: time.Minute,
				Host:    "https://tmdb.host",
				Token:   "",
			}},
			want: "Key: 'Config.Token' Error:Field validation for 'Token' failed on the 'required' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := tmdb.New(test.args.config)

			if test.want != "" {
				failure(t, obj, err, test.want)
			} else {
				success(t, obj, err)
			}
		})
	}
}

func TestTMDBSetTransport(t *testing.T) {
	t.Parallel()

	obj1 := New()
	obj2 := obj1.SetTransport(new(http.Transport))
	obj3 := obj1.SetTransport(mocks.NewMockRoundTripper(t))

	assert.Same(t, obj1, obj2)
	assert.Same(t, obj1, obj3)
}

func TestTMDBClose(t *testing.T) {
	t.Parallel()

	err := New().Close()

	require.NoError(t, err)
}

func success(t *testing.T, got *tmdb.TMDB, err error) {
	t.Helper()

	require.NoError(t, err)

	assert.IsType(t, new(tmdb.TMDB), got)
}

func failure(t *testing.T, obj *tmdb.TMDB, got error, want string) {
	t.Helper()

	var orr oops.OopsError

	require.ErrorAs(t, got, &orr)
	require.EqualError(t, got, want)

	assert.Equal(t, "Invalid configuration.", orr.Public())
	assert.Nil(t, obj)
}
