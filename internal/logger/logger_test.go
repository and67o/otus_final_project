package logger

import (
	"testing"

	"github.com/and67o/otus_project/internal/configuration"
	"github.com/stretchr/testify/require"
)

func TestLogger(t *testing.T) {
	t.Run("Ok tests", func(t *testing.T) {
		_, err := New(configuration.LoggerConf{
			Level:   levelDebug,
			File:    "./testdata/log.log",
			IsProd:  false,
			TraceOn: false,
		})
		require.NoError(t, err)
	})

	t.Run("fail path", func(t *testing.T) {
		_, err := New(configuration.LoggerConf{
			Level:   "",
			File:    "",
			IsProd:  false,
			TraceOn: false,
		})
		require.Error(t, err)
	})
}
