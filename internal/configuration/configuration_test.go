package configuration

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	path string
	err  error
}

func TestConfiguration(t *testing.T) {
	t.Run("errors", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				path: "",
				err:  errEmptyPath,
			},
		} {
			_, err := New(tst.path)
			require.Equal(t, tst.err, err)
		}
	})

	t.Run("wrong path", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				path: "./wrong_path/config_test.toml",
			},
			{
				path: "wrong_path",
			},
		} {
			_, err := New(tst.path)
			require.Error(t, err)
		}
	})

	t.Run("pass result", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				path: "./testdata/config_test.toml",
			},
		} {
			c, err := New(tst.path)
			require.Nil(t, err)

			require.Equal(t, c.DB.User, "db")
			require.Equal(t, c.Logger.Level, "log")
			require.Equal(t, c.Server.Host, "server")

			require.Equal(t, c.Rabbit.User, "")
		}
	})
}
