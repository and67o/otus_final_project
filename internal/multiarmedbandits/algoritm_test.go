package multiarmedbandits

import (
	"testing"

	"github.com/and67o/otus_project/internal/model"
	"github.com/stretchr/testify/require"
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
			input: []model.Banner{
				{ID: 1, SlotID: 2, ShowCount: 3, ClickCount: 4},
			},
			expected: 1,
		},
		{
			input: []model.Banner{
				{ID: 1, SlotID: 2, ShowCount: 3, ClickCount: 4},
				{ID: 1, SlotID: 2, ShowCount: 4, ClickCount: 8},
				{ID: 2, SlotID: 2, ShowCount: 34, ClickCount: 41},
				{ID: 3, SlotID: 2, ShowCount: 99, ClickCount: 123},
				{ID: 4, SlotID: 2, ShowCount: 1, ClickCount: 1},
				{ID: 8, SlotID: 2, ShowCount: 143, ClickCount: 156},
				{ID: 5, SlotID: 2, ShowCount: 23, ClickCount: 45},
				{ID: 6, SlotID: 2, ShowCount: 23, ClickCount: 1},
				{ID: 7, SlotID: 2, ShowCount: 67, ClickCount: 67},
				{ID: 9, SlotID: 2, ShowCount: 122, ClickCount: 34},
			},
			expected: 4,
		},
	} {
		require.Equal(t, tst.expected, Get(tst.input))
	}
}
