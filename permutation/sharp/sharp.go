package sharp

import (
	"fmt"
)

const (
	KeyLength = 4
)

type Point struct {
	x, y int
}

type SharpCrypto struct {
	key map[int]Point
}

func New(key [KeyLength]Point) *SharpCrypto {
	var (
		crypto = &SharpCrypto{
			key: make(map[int]Point),
		}

		matrix = [4][4]int{
			{0, 1, 2, 0},
			{2, 3, 3, 1},
			{1, 3, 3, 2},
			{0, 2, 1, 0},
		}
	)

	for _, point := range key {
		crypto.key[matrix[point.x][point.y]] = point
	}

	if len(crypto.key) != KeyLength {
		panic("invalid point in the key")
	}

	return crypto
}

func (s *SharpCrypto) Encode(dst, src []rune) {
	result := make([]rune, 0, len(dst))

	s.split(src, 16, func(r []rune) {
		res := s.encode(r)

		result = append(result, res...)
	})

	copy(dst, result)
}

func (s *SharpCrypto) encode(src []rune) []rune {
	if len(src) != 16 {
		panic(fmt.Sprintf("invalid split 16 str: %d", len(src)))
	}

	var (
		count  = 0
		matrix = NewEmptyMatrix()
	)

	s.split(src, KeyLength, func(r []rune) {

		for key, value := range s.key {
			matrix.Add(subTable(value, 4-count), key, r[key])
		}

		count++
	})

	result := make([]rune, 0, len(src))

	for _, str := range matrix.Result() {
		result = append(result, str...)
	}

	return result
}

func subTable(point Point, turnCount int) int {
	if point.y < 2 {
		if point.x < 2 {
			return turnCount % 4
		}

		return (3 + turnCount) % 4
	}

	if point.x < 2 {
		return (1 + turnCount) % 4
	}

	return (2 + turnCount) % 4
}

func (k *SharpCrypto) split(src []rune, length int, f func(r []rune)) {
	i := 0
	for ; i < len(src)/length; i++ {
		f(src[i*length : (i+1)*length])
	}

	end := src[i*length:]
	if len(end) > 0 {
		if len(end) != length {
			end = append(end, k.generateEnd(length-len(end))...)
		}

		f(end)
	}

}

func (k *SharpCrypto) generateEnd(length int) []rune {
	const placeholder = '_'

	result := make([]rune, 0, length)

	for i := 0; i < length; i++ {
		result = append(result, placeholder)
	}

	return result
}

func (k *SharpCrypto) Decode(dst, src []rune) {
	result := make([]rune, 0, len(src))

	k.split(src, 16, func(r []rune) {
		res := k.decodeBlock(r)

		result = append(result, res...)
	})

	copy(dst, result)
}

func (k *SharpCrypto) decodeBlock(src []rune) []rune {
	if len(src) != 16 {
		panic(fmt.Sprintf("invalid block size[%d]", len(src)))
	}

	var (
		matrix = k.createMatrix(src)
		result = make([]rune, 0, len(src))
	)

	for count := 0; count < 4; count++ {
		res := make([]rune, 4)

		for key, value := range k.key {
			res[key] = matrix.Get(subTable(value, 4-count), key)
		}

		result = append(result, res...)
	}

	return result
}

func (k *SharpCrypto) createMatrix(src []rune) *Matrix {
	var (
		count = 0
		table = make([][]rune, 4)
	)

	k.split(src, 4, func(r []rune) {
		table[count] = append(table[count], r...)
		count++
	})

	return NewFullMatrix(table)
}
