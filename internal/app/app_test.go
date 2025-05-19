package app_test

import (
	"bufio"
	"io"
	"net/http"
	"testing"
	"time"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/tmdb/internal/app"
	"github.com/therenotomorrow/tmdb/internal/app/mocks"
	"github.com/therenotomorrow/tmdb/internal/config"
	"github.com/therenotomorrow/tmdb/pkg/fp"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

func New() *app.TMDB {
	return fp.Must(app.New(config.Settings{
		Debug: false,
		Token: "secret",
		Config: tmdb.Config{
			Host:    "https://tmdb.host",
			Token:   "secret",
			Timeout: time.Minute,
			Debug:   false,
		},
	}))
}

func TestNew(t *testing.T) {
	t.Parallel()

	type args struct {
		settings config.Settings
	}

	tests := []struct {
		name string
		want string
		args args
	}{
		{
			name: "min settings",
			args: args{
				settings: config.Settings{
					Token: "secret",
					Config: tmdb.Config{
						Host:    "https://tmdb.host",
						Token:   "secret",
						Timeout: time.Minute,
						Debug:   false,
					},
					Debug: false,
				},
			},
			want: "",
		},
		{
			name: "max settings",
			args: args{
				settings: config.Settings{
					Token: "secret",
					Config: tmdb.Config{
						Host:    "https://tmdb.host",
						Token:   "secret",
						Timeout: time.Minute,
						Debug:   true,
					},
					Debug: true,
				},
			},
			want: "",
		},
		{
			name: "settings error",
			args: args{
				settings: config.Settings{
					Token: "",
					Config: tmdb.Config{
						Host:    "https://tmdb.host",
						Token:   "",
						Timeout: time.Minute,
						Debug:   true,
					},
					Debug: true,
				},
			},
			want: "Key: 'Config.Token' Error:Field validation for 'Token' failed on the 'required' tag",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			obj, err := app.New(test.args.settings)

			if test.want != "" {
				failure(t, obj, err, test.want)
			} else {
				success(t, obj, err)
			}
		})
	}
}

func success(t *testing.T, got *app.TMDB, err error) {
	t.Helper()

	require.NoError(t, err)

	assert.IsType(t, new(app.TMDB), got)
}

func failure(t *testing.T, obj *app.TMDB, got error, want string) {
	t.Helper()

	var orr oops.OopsError

	require.ErrorAs(t, got, &orr)
	require.EqualError(t, got, want)

	assert.Equal(t, "Invalid configuration.", orr.Public())
	assert.Nil(t, obj)
}

func TestTMDBWithDependencies(t *testing.T) {
	t.Parallel()

	obj1 := New()
	obj2 := obj1.WithDependencies(mocks.NewMockClient(t), bufio.NewReader(nil), bufio.NewWriter(nil))
	obj3 := obj1.WithDependencies(new(http.Transport), io.Reader(nil), io.Writer(nil))

	assert.Same(t, obj1, obj2)
	assert.Same(t, obj1, obj3)
}

func TestTMDBClose(t *testing.T) {
	t.Parallel()

	err := New().Close()

	require.NoError(t, err)
}
