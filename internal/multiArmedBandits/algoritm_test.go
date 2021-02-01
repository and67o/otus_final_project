package multiArmedBandits

import (
	"github.com/and67o/otus_project/internal/model"
	"github.com/stretchr/testify/require"
	"testing"
)

type test struct {
	input    []model.Banner
	expected int64
}

func TestGet(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    []model.Banner{},
			expected: 0,
		},
		{
			input:    []model.Banner{{1, 2, 3, 4}},
			expected: 1,
		},
		{
			input: []model.Banner{{1, 2, 3, 4},
				{1, 2, 4, 8},
				{2, 2, 34, 41},
				{3, 2, 99, 123},
				{4, 2, 1, 1},
				{8, 2, 143, 156},
				{5, 2, 23, 45},
				{6, 2, 23, 1},
				{7, 2, 67, 67},
				{9, 2, 122, 34},
			},
			expected: 4,
		},
	} {
		require.Equal(t, tst.expected, Get(tst.input))
	}
}
