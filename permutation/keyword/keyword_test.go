package keyword

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOrder(t *testing.T) {
	tt := []struct {
		key  string
		want map[int]int
	}{
		{
			key: "криптография",
			want: map[int]int{
				0:  8,
				1:  6,
				2:  2,
				3:  10,
				4:  0,
				5:  5,
				6:  3,
				7:  1,
				8:  7,
				9:  4,
				10: 9,
				11: 11,
			},
		},
	}

	for _, tc := range tt {
		t.Run(tc.key, func(t *testing.T) {
			k := New(tc.key)

			assert.Equal(t, tc.want, k.order, "order not equel")
		})
	}
}
