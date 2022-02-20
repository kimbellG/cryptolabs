package sharp

import (
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestSubTable(t *testing.T) {
	tt := []struct {
		count int
		point Point

		want int
	}{
		{
			count: 1,
			point: Point{x: 0, y: 0},

			want: 1,
		},
		{
			count: 1,
			point: Point{x: 2, y: 2},

			want: 3,
		},
		{
			count: 2,
			point: Point{x: 2, y: 2},

			want: 0,
		},
		{
			count: 0,
			point: Point{x: 1, y: 3},

			want: 1,
		},
	}

	for _, tc := range tt {
		got := subTable(tc.point, tc.count)

		assert.Equal(t, tc.want, got, "invalid subtable from matrix")
	}
}

func TestEncode(t *testing.T) {
	tt := []struct {
		key  [4]Point
		arg  string
		want string
	}{
		{
			key: [4]Point{
				{x: 0, y: 0},
				{x: 1, y: 3},
				{x: 2, y: 2},
				{x: 3, y: 1},
			},
			arg:  "ЭТОЛЕКЦИЯПОКРИПТ",
			want: "ЭКОРПКИТПТЛЦЕОИЯ",
		},
		{
			key: [4]Point{
				{0, 3},
				{1, 1},
				{2, 0},
				{3, 1},
			},

			arg:  "АРТËМФИЛИППЕН",
			want: "М_ПА_Ë_ПРЛЕИИТФН",
		},
	}

	for _, tc := range tt {
		var (
			sharp = New(tc.key)
			got   = make([]rune, utf8.RuneCountInString(tc.want))
		)

		sharp.Encode(got, []rune(tc.arg))

		assert.Equal(t, tc.want, string(got), "invalid encryption")
	}
}

func TestDecode(t *testing.T) {
	tt := []struct {
		key    [4]Point
		arg    string
		encLen int
	}{
		{
			key: [4]Point{
				{x: 0, y: 0},
				{x: 1, y: 3},
				{x: 2, y: 2},
				{x: 3, y: 1},
			},
			encLen: 16,
			arg:    "ЭТОЛЕКЦИЯПОКРИПТ",
		},
		{
			key: [4]Point{
				{0, 3},
				{1, 1},
				{2, 0},
				{3, 1},
			},

			encLen: 16,
			arg:    "АРТËМФИЛИППЕН",
		},
	}

	for _, tc := range tt {
		var (
			sharp = New(tc.key)
			enc   = make([]rune, tc.encLen)
			got   = make([]rune, tc.encLen)
		)

		sharp.Encode(enc, []rune(tc.arg))
		sharp.Decode(got, enc)

		got = []rune(strings.TrimRight(string(got), "_"))
		assert.Equalf(t, tc.arg, string(got), "invalid encryption for %s", tc.arg)
	}

}
