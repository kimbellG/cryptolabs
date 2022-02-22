package extcaesar

import (
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
)

func TestEncode(t *testing.T) {
	tt := []struct {
		enc int
		dec int

		arg  string
		want string
	}{
		{
			enc:  7,
			dec:  15,
			arg:  "artyom",
			want: "apdmug",
		},
	}

	for _, tc := range tt {
		var (
			caesar = New(tc.enc, tc.dec)
			got    = make([]rune, utf8.RuneCountInString(tc.arg))
		)

		caesar.Encode(got, []rune(tc.arg))

		assert.Equal(t, tc.want, string(got), "invalid caesar encryption")
	}
}

func TestDecode(t *testing.T) {
	tt := []struct {
		enc, dec int

		arg string
	}{
		{
			enc: 7,
			dec: 15,
			arg: "cryptographyanddatasecurity",
		},
	}

	for _, tc := range tt {
		var (
			caesar = New(tc.enc, tc.dec)
			enc    = make([]rune, utf8.RuneCountInString(tc.arg))
			got    = make([]rune, utf8.RuneCountInString(tc.arg))
		)

		caesar.Encode(enc, []rune(tc.arg))
		caesar.Decode(got, enc)

		assert.Equal(t, tc.arg, string(got), "invalid caesar encryption")
	}
}
