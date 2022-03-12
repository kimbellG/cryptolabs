package murmurhash3

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSum32(t *testing.T) {
	tt := []struct {
		input []byte
		seed  uint32

		want uint32
	}{
		{
			input: []byte{'B', 'S', 'U', 'I', 'R'},
			seed:  21,

			want: 0xCF31D00D,
		},
	}

	for _, tc := range tt {
		got := Sum32WithSeed(tc.input, tc.seed)

		assert.Equal(t, tc.want, got, "")
	}
}
