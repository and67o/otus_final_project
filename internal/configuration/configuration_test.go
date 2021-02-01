package configuration

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	path string
	err  error
}

func TestConfiguration(t *testing.T) {
	t.Run("Zero N and M", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				path: "",
				err:  errors.New("path empty"),
			},
		} {
			_, err := New(tst.path)
			require.Equal(t, tst.err, err)
		}
	})

	t.Run("Zero N and M", func(t *testing.T) {
		for _, tst := range [...]test{
			{
				path: "./test/config_test.toml",
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
