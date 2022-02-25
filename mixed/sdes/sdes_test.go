package sdes

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRearrange(t *testing.T) {
	tt := []struct {
		order  []int
		src    []byte
		srcLen int
		want   []byte
	}{
		{
			order:  []int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6},
			srcLen: 10,
			src:    []byte{83, 2},
			want:   []byte{61, 0},
		},
		{

			order:  []int{6, 3, 7, 4, 8, 5, 10, 9},
			srcLen: 10,
			src:    []byte{0x5B, 0},
			want:   []byte{0xB3},
		},
	}

	for _, tc := range tt {
		got := rearrange(tc.order, tc.src, tc.srcLen)

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

func TestContact5Bit(t *testing.T) {
	tt := []struct {
		l, r byte

		want []byte
	}{
		{
			l:    0x2,
			r:    0x1B,
			want: []byte{0x5B, 0x0},
		},
	}

	for _, tc := range tt {
		got := contact5Bit(tc.l, tc.r)

		assert.Equal(t, tc.want, got, "")
	}
}

func TestRotate(t *testing.T) {
	tt := []struct {
		src  byte
		k    int
		want byte
	}{
		{
			src:  1,
			k:    1,
			want: 2,
		},
		{
			src:  1,
			k:    -1,
			want: 0x10,
		},
		{
			src:  0x1D,
			k:    6,
			want: 0x1B,
		},
	}

	for _, tc := range tt {
		got := rotateLeft5(tc.src, tc.k)

		assert.Equal(t, tc.want, got, "invalid rotate")
	}

}

func TestKeyGenerate(t *testing.T) {
	tt := []struct {
		key    []byte
		k1, k2 byte
	}{
		{
			key: []byte{0x53, 0x2},
			k1:  0xB3,
			k2:  0x2B,
		},
	}

	for _, tc := range tt {
		k1, k2 := generateKeys(tc.key)

		assert.Equal(t, tc.k1, k1, "invalid the first key")
		assert.Equal(t, tc.k2, k2, "invalid the second key")

	}
}
