package sharp

import "fmt"

type Matrix struct {
	indexes map[int][][]int
	values  map[int][][]rune
}

func newIndexes() map[int][][]int {
	matrix := make(map[int][][]int)

	matrix[0] = [][]int{
		{0, 1},
		{2, 3},
	}

	matrix[1] = [][]int{
		{2, 0},
		{3, 1},
	}

	matrix[3] = [][]int{
		{1, 3},
		{0, 2},
	}

	matrix[2] = [][]int{
		{3, 2},
		{1, 0},
	}

	return matrix
}

func NewEmptyMatrix() *Matrix {

	values := make(map[int][][]rune)

	for i := 0; i < KeyLength; i++ {
		values[i] = make([][]rune, 2)

		for j := 0; j < 2; j++ {
			values[i][j] = make([]rune, 2)
		}
	}

	return &Matrix{
		indexes: newIndexes(),
		values:  values,
	}
}

func NewFullMatrix(r [][]rune) *Matrix {
	if len(r) != 4 {
		panic(fmt.Sprintf("new full matrix: invalid length of the matrix: %d", len(r)))
	}

	m := NewEmptyMatrix()

	copy(m.values[0][0], r[0][:3])
	copy(m.values[1][0], r[0][2:])

	copy(m.values[0][1], r[1][:3])
	copy(m.values[1][1], r[1][2:])

	copy(m.values[3][0], r[2][:3])
	copy(m.values[2][0], r[2][2:])

	copy(m.values[3][1], r[3][:3])
	copy(m.values[2][1], r[3][2:])

	return m
}

func (m *Matrix) Add(table, index int, r rune) {
	for i := 0; i < len(m.indexes[table]); i++ {
		for j := 0; j < len(m.indexes[table][i]); j++ {
			if m.indexes[table][i][j] == index {
				m.values[table][i][j] = r
			}
		}
	}
}

func (m *Matrix) Get(table, index int) rune {
	for i := 0; i < len(m.indexes[table]); i++ {
		for j := 0; j < len(m.indexes[table][i]); j++ {
			if m.indexes[table][i][j] == index {
				return m.values[table][i][j]
			}
		}
	}

	panic(fmt.Sprintf("invalid table[%d] or index[%d]", table, index))
}

func (m *Matrix) Result() [][]rune {
	result := make([][]rune, KeyLength)

	result[0] = append(result[0], m.values[0][0]...)
	result[0] = append(result[0], m.values[1][0]...)

	result[1] = append(result[1], m.values[0][1]...)
	result[1] = append(result[1], m.values[1][1]...)

	result[2] = append(result[2], m.values[3][0]...)
	result[2] = append(result[2], m.values[2][0]...)

	result[3] = append(result[3], m.values[3][1]...)
	result[3] = append(result[3], m.values[2][1]...)

	return result
}
