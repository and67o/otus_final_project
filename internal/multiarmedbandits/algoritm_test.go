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
				{IDBanner: 1, IDSlot: 2, CountShow: 3, CountClick: 4},
			},
			expected: 1,
		},
		{
			input: []model.Banner{
				{IDBanner: 1, IDSlot: 2, CountShow: 3, CountClick: 4},
				{IDBanner: 1, IDSlot: 2, CountShow: 4, CountClick: 8},
				{IDBanner: 2, IDSlot: 2, CountShow: 34, CountClick: 41},
				{IDBanner: 3, IDSlot: 2, CountShow: 99, CountClick: 123},
				{IDBanner: 4, IDSlot: 2, CountShow: 1, CountClick: 1},
				{IDBanner: 8, IDSlot: 2, CountShow: 143, CountClick: 156},
				{IDBanner: 5, IDSlot: 2, CountShow: 23, CountClick: 45},
				{IDBanner: 6, IDSlot: 2, CountShow: 23, CountClick: 1},
				{IDBanner: 7, IDSlot: 2, CountShow: 67, CountClick: 67},
				{IDBanner: 9, IDSlot: 2, CountShow: 122, CountClick: 34},
			},
			expected: 4,
		},
	} {
		require.Equal(t, tst.expected, Get(tst.input))
	}
}
