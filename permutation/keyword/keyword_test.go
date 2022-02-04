package keyword

import (
	"testing"
	"unicode/utf8"

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

func TestEncode(t *testing.T) {
	tt := []struct {
		key  string
		arg  string
		want string
	}{
		{
			key:  "КРИПТОГРАФИЯ",
			arg:  "ЭТО–_ЛЕКЦИЯ_ПО_АЛГОРИТМАМ_ШИФРОВАНИЯ",
			want: "ЦЕОЯЭЛ–ТК_И_ИО_МПГАОРЛТААОШИМРИ_ВФНЯ",
		},
	}

	for _, tc := range tt {
		var (
			crypto = New(tc.key)
			got    = make([]rune, len([]rune(tc.arg)))
		)

		crypto.Encode(got, []rune(tc.arg))

		assert.Equal(t, tc.want, string(got), "incorrect encryption")
	}
}

func TestDecode(t *testing.T) {
	tt := []struct {
		key string
		arg string
	}{
		{
			key: "Артём",
			arg: "ЭТО–_ЛЕКЦИЯ_ПО_АЛГОРИТМАМ_ШИФРОВАНИЯ",
		},
	}

	for _, tc := range tt {
		var (
			crypto     = New(tc.key)
			got        = make([]rune, utf8.RuneCountInString(tc.arg))
			encryption = make([]rune, utf8.RuneCountInString(tc.arg))
		)

		crypto.Encode(encryption, []rune(tc.arg))
		crypto.Decode(got, encryption)

		assert.Equal(t, tc.arg, string(got), "decoding is incorrect")
	}
}
