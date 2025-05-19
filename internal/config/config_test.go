package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/samber/oops"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/tmdb/internal/config"
	"github.com/therenotomorrow/tmdb/pkg/tmdb"
)

func want() config.Settings {
	return config.Settings{
		Debug: true,
		Token: "secret",
		Config: tmdb.Config{
			Debug:   true,
			Timeout: 10 * time.Second,
			Host:    "https://api.themoviedb.org",
			Token:   "secret",
		},
	}
}

func TestNewSuccessFromDotenv(t *testing.T) {
	t.Setenv(t.Name(), "test")

	const filename = ".env.test"

	_ = os.WriteFile(filename, []byte("TMDB_DEBUG=true\nTMDB_TOKEN=secret\n"), 0o600)
	defer func() { _ = os.Remove(filename) }()

	got, err := config.New(filename)

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestNewSuccessFromEnvironment(t *testing.T) {
	t.Setenv(t.Name(), "test")

	t.Setenv("TMDB_DEBUG", "true")
	t.Setenv("TMDB_TOKEN", "secret")

	got, err := config.New("skip")

	require.NoError(t, err)
	assert.Equal(t, want(), got)
}

func TestNewFailure(t *testing.T) {
	t.Setenv(t.Name(), "test")

	t.Setenv("TMDB_DEBUG", "invalid")
	t.Setenv("TMDB_TOKEN", "")

	var orr oops.OopsError

	got, err := config.New("skip")

	require.ErrorAs(t, err, &orr)
	require.EqualError(t, err, `Debug: strconv.ParseBool: parsing "invalid": invalid syntax`)

	assert.Equal(t, "Invalid configuration.", orr.Public())
	assert.Empty(t, got)
}

func TestSettingsSetToken(t *testing.T) {
	t.Parallel()

	const want = "supersecret"

	var obj config.Settings

	obj.SetToken(want)

	assert.Equal(t, want, obj.Token)
	assert.Equal(t, want, obj.Config.Token)
}
