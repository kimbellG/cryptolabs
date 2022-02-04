package keyword

import (
	"unicode/utf8"
)

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
	result := make([]rune, 0, len(src))

	k.split(src, func(r []rune) {
		if k.isEnd(r) {
			result = append(result, r...)

			return
		}

		res := make([]rune, k.keyLength)
		for key, value := range k.order {
			res[key] = r[value]
		}

		result = append(result, res...)
	})

	copy(dst, result)
}

func (k *KeywordCrypto) split(src []rune, f func(r []rune)) {
	i := 0
	for ; i < len(src)/k.keyLength; i++ {
		f(src[i*k.keyLength : (i+1)*k.keyLength])
	}

	f(src[i*k.keyLength:])
}

func (k *KeywordCrypto) isEnd(r []rune) bool {
	return len(r) < k.keyLength
}

func (k *KeywordCrypto) Decode(dst, src []rune) {
	result := make([]rune, 0, len(src))

	k.split(src, func(r []rune) {
		if k.isEnd(r) {
			result = append(result, r...)

			return
		}

		res := make([]rune, k.keyLength)
		for key, value := range k.order {
			res[value] = r[key]
		}

		result = append(result, res...)
	})

	copy(dst, result)
}
