package extcaesar

import (
	"fmt"
)

const AlfaLength = 26

type Caesar struct {
	enc   int
	dec   int
	index map[rune]int
	runes map[int]rune
}

func New(enc, dec int) *Caesar {
	return &Caesar{
		enc: enc,
		dec: dec,
		index: map[rune]int{
			'a': 0,
			'b': 1,
			'c': 2,
			'd': 3,
			'e': 4,
			'f': 5,
			'g': 6,
			'h': 7,
			'i': 8,
			'j': 9,
			'k': 10,
			'l': 11,
			'm': 12,
			'n': 13,
			'o': 14,
			'p': 15,
			'q': 16,
			'r': 17,
			's': 18,
			't': 19,
			'u': 20,
			'v': 21,
			'w': 22,
			'x': 23,
			'y': 24,
			'z': 25,
		},
		runes: map[int]rune{
			0:  'a',
			1:  'b',
			2:  'c',
			3:  'd',
			4:  'e',
			5:  'f',
			6:  'g',
			7:  'h',
			8:  'i',
			9:  'j',
			10: 'k',
			11: 'l',
			12: 'm',
			13: 'n',
			14: 'o',
			15: 'p',
			16: 'q',
			17: 'r',
			18: 's',
			19: 't',
			20: 'u',
			21: 'v',
			22: 'w',
			23: 'x',
			24: 'y',
			25: 'z',
		},
	}
}

func (c *Caesar) Encode(dst, src []rune) {
	c.validate(src)

	for i := range src {
		dst[i] = c.encodeChar(src[i])
	}
}

func (c *Caesar) encodeChar(r rune) rune {
	return c.runes[(c.index[r]*c.enc)%AlfaLength]
}

func (c *Caesar) validate(src []rune) {
	for _, r := range src {
		if _, ok := c.index[r]; !ok {
			panic(fmt.Sprintf("caesar: invalid character in the source string[%c]", r))
		}
	}
}

func (c *Caesar) Decode(dst, src []rune) {
	c.validate(src)

	for i := range src {
		dst[i] = c.decodeChar(src[i])
	}
}

func (c *Caesar) decodeChar(r rune) rune {
	return c.runes[(c.index[r]*c.dec)%AlfaLength]
}
