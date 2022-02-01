package rail

import (
	"testing"
	"unicode/utf8"
)

func TestRailEncoder(t *testing.T) {
	tt := []struct {
		key  int
		arg  string
		want string
	}{
		{
			key:  3,
			arg:  "Артём Филиппенков",
			want: "Амлеврё иипнотФпк",
		},
		{
			key:  4,
			arg:  "Артём Филиппенков",
			want: "АФер ипнтмлпквёио",
		},
	}

	for _, tc := range tt {
		got := make([]rune, utf8.RuneCountInString(tc.arg))

		encoder := New(tc.key)
		encoder.Encode(got, []rune(tc.arg))

		if string(got) != tc.want {
			t.Errorf("encryption failed: want: \"%s\". got: \"%s\"", tc.want, string(got))
		}
	}
}

func TestRailDecoder(t *testing.T) {
	tt := []struct {
		key int
		arg string
	}{
		{
			key: 3,
			arg: "Артём Филиппенков",
		},
		{
			key: 4,
			arg: "Артём Филиппенков",
		},
	}

	for _, tc := range tt {
		enc := make([]rune, utf8.RuneCountInString(tc.arg))
		got := make([]rune, utf8.RuneCountInString(tc.arg))

		encoder := New(tc.key)
		encoder.Encode(enc, []rune(tc.arg))
		encoder.Decode(got, enc)

		if string(got) != tc.arg {
			t.Errorf("encryption failed: want: \"%s\". got: \"%s\"", tc.arg, string(got))
		}
	}
}
