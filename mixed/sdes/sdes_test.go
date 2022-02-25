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
		{
			order:  []int{4, 1, 2, 3, 2, 3, 4, 1},
			srcLen: 4,
			src:    []byte{0x7},
			want:   []byte{0xBE},
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

func TestHalve8Bits(t *testing.T) {
	tt := []struct {
		b    byte
		l, r byte
	}{
		{
			b: 0xBF,
			l: 0xB,
			r: 0xF,
		},
	}

	for _, tc := range tt {
		l, r := halve8Bit(tc.b)

		assert.Equal(t, tc.l, l, "invalid left bits")
		assert.Equal(t, tc.r, r, "invalid right bits")
	}
}

func TestContact4Bits(t *testing.T) {
	tt := []struct {
		left, right byte
		want        byte
	}{
		{
			left:  0xF,
			right: 0x3,
			want:  0xF3,
		},
	}

	for _, tc := range tt {
		got := contact4Bit(tc.left, tc.right)

		assert.Equal(t, tc.want, got, "")
	}
}

func TestGetSBlockCord(t *testing.T) {
	tt := []struct {
		b        byte
		row, col byte
	}{
		{
			b:   0x9,
			row: 3,
			col: 0,
		},
	}

	for _, tc := range tt {
		row, col := getSBlockCord(tc.b)

		assert.Equal(t, tc.row, row, "invalid row")
		assert.Equal(t, tc.col, col, "invalid col")
	}
}

func TestSBlock(t *testing.T) {
	tt := []struct {
		sblock [][]byte
		b      byte
		want   byte
	}{
		{
			sblock: [][]byte{
				{1, 0, 3, 2},
				{3, 2, 1, 0},
				{0, 2, 1, 3},
				{3, 1, 3, 2},
			},
			b:    0x9,
			want: 0x3,
		},
	}

	for _, tc := range tt {
		got := blockproc(tc.sblock, tc.b)

		assert.Equal(t, tc.want, got, "")
	}
}

func TestRound(t *testing.T) {
	tt := []struct {
		input byte
		key   byte

		want byte
	}{
		{
			input: 0x79,
			key:   0xB3,

			want: 0x79,
		},
	}

	for _, tc := range tt {
		block := SDES{}

		block.s1 = [][]byte{
			{1, 0, 3, 2},
			{3, 2, 1, 0},
			{0, 2, 1, 3},
			{3, 1, 3, 2},
		}
		block.s2 = [][]byte{
			{0, 1, 2, 3},
			{2, 0, 1, 3},
			{3, 0, 1, 0},
			{2, 1, 0, 3},
		}

		got := block.round(tc.input, tc.key)

		assert.Equal(t, tc.want, got, "")
	}
}

func TestEncrypt(t *testing.T) {
	tt := []struct {
		key   []byte
		input byte

		want byte
	}{
		{
			key:   []byte{0x53, 0x2},
			input: 0xB6,

			want: 0x0F,
		},
	}

	for _, tc := range tt {
		block := New(tc.key)

		got := block.encrypt(tc.input)
		assert.Equal(t, tc.want, got, "")
	}
}

func TestDecrypt(t *testing.T) {
	tt := []struct {
		key   []byte
		input byte

		want byte
	}{
		{
			key:   []byte{0x53, 0x2},
			input: 0x0F,

			want: 0xB6,
		},
	}

	for _, tc := range tt {
		block := New(tc.key)

		got := block.decrypt(tc.input)
		assert.Equal(t, tc.want, got, "")
	}
}
