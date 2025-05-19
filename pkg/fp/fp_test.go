package fp_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/therenotomorrow/tmdb/pkg/fp"
)

var errOops = errors.New("oops")

func TestMustSuccess(t *testing.T) {
	t.Parallel()

	var got int

	require.NotPanics(t, func() {
		got = fp.Must(func() (int, error) { return 42, nil }())
	})

	want := 42

	assert.Equal(t, want, got)
}

func TestMustFailure(t *testing.T) {
	t.Parallel()

	var got int

	require.Panics(t, func() {
		got = fp.Must(func() (int, error) { return 0, errOops }())
	})

	want := 0

	assert.Equal(t, want, got)
}

func TestSilent(t *testing.T) {
	t.Parallel()

	fp.Silent(func() (int, error) { return 42, errOops }())
}
