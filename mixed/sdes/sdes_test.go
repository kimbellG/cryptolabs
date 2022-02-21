package sdes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRearrange(t *testing.T) {
	tt := []struct {
		order []int
		src   []byte
		want  []byte
	}{
		{
			order: []int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6},
			src:   []byte{83, 2},
			want:  []byte{61, 0},
		},
	}

	for _, tc := range tt {
		got := rearrange(tc.order, tc.src)

		assert.Equal(t, tc.want, got, "invalid rearrange bits")
	}
}

func TestHalve10Bits(t *testing.T) {
	tt := []struct {
		src []byte

		leftWant, rightWant byte
	}{
		{
			src:       []byte{0x3D, 0x0},
			leftWant:  1,
			rightWant: 0x1D,
		},
	}

	for _, tc := range tt {
		left, right := halve10Bit(tc.src)

		assert.Equal(t, tc.leftWant, left, "invalid left part")
		assert.Equal(t, tc.rightWant, right, "invalid left part")
	}

}
