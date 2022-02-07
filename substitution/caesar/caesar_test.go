package caesar

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	tt := []struct {
		key int

		arg  string
		want string
	}{
		{
			key:  3,
			arg:  "cryptographyanddatasecurity",
			want: "fubswrjudskbdqggdwdvhfxulwb",
		},
	}

	for _, tc := range tt {
		var (
			caesar = New(tc.key)
			got    = make([]rune, utf8.RuneCountInString(tc.arg))
		)

		caesar.Encode(got, []rune(tc.arg))

		assert.Equal(t, tc.want, string(got), "invalid caesar encryption")
	}
}

func TestDecode(t *testing.T) {
	tt := []struct {
		key int

		arg string
	}{
		{
			key: 3,
			arg: "cryptographyanddatasecurity",
		},
	}

	for _, tc := range tt {
		var (
			caesar = New(tc.key)
			enc    = make([]rune, utf8.RuneCountInString(tc.arg))
			got    = make([]rune, utf8.RuneCountInString(tc.arg))
		)

		caesar.Encode(enc, []rune(tc.arg))
		caesar.Decode(got, enc)

		assert.Equal(t, tc.arg, string(got), "invalid caesar encryption")
	}
}
