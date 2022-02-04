package sharp

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
			{1, 2, 3, 1},
			{3, 4, 4, 2},
			{2, 4, 4, 3},
			{1, 3, 2, 1},
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

}
