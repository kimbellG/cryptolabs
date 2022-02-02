package keyword

import "unicode/utf8"

type KeywordCrypto struct {
	order     map[int]int
	keyLength int
}

func New(key string) *KeywordCrypto {
	keyw := &KeywordCrypto{
		keyLength: utf8.RuneCountInString(key),
	}

	keyw.order = keyw.toOrder([]rune(key))

	return keyw
}

func (k *KeywordCrypto) toOrder(key []rune) map[int]int {
	result := make(map[int]int)

	for i := range []rune(key) {
		lower := k.lower(key)

		result[i] = lower
		key[lower] = '\uFFFD'
	}

	return result
}

func (k *KeywordCrypto) lower(key []rune) int {
	var (
		index = 0
		lower = key[index]
	)

	for i, ch := range key {
		if ch < lower {
			lower, index = ch, i
		}
	}

	return index
}

func (k *KeywordCrypto) Encode(dst, src []rune) {
	result := make([]rune, 0, k.keyLength)

	

}

func (k *KeywordCrypto) split(src []rune, f func(r []rune)) {
	for i := 0; i < len(src) / k.keyLength - 1; i += k.keyLength {
		f(src
	}
}
