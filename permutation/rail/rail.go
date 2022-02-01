package rail

import "log"

type railTable struct {
	buf    [][]rune
	height int
	key    int

	isHeightDec bool
}

func newTable(key int) *railTable {
	return &railTable{
		buf: make([][]rune, key),
		key: key,
	}
}

func (t *railTable) addSlice(rs []rune) {
	for _, r := range rs {
		t.add(r)
	}
}

func (t *railTable) add(r rune) {
	t.buf[t.height] = append(t.buf[t.height], r)

	t.next()
}

func (t *railTable) next() {
	if t.isDown() {
		t.isHeightDec = true
	}

	if t.isUp() {
		t.isHeightDec = false
	}

	if t.isHeightDec {
		t.height--
	} else {
		t.height++
	}
}

func (t *railTable) table() [][]rune {
	return t.buf
}

func (t *railTable) isDown() bool {
	return t.height == t.key-1
}

func (t *railTable) isUp() bool {
	return t.height == 0
}

type RailCrypto struct {
	key int
}

func New(key int) *RailCrypto {
	return &RailCrypto{
		key: key,
	}
}

func (r *RailCrypto) Encode(dst, src []rune) {
	railTable := r.createRail(src)

	r.writeEncryption(railTable, dst)
}

func (r *RailCrypto) createRail(src []rune) [][]rune {
	table := newTable(r.key)
	table.addSlice(src)

	return table.table()
}

func (r *RailCrypto) writeEncryption(railTable [][]rune, dst []rune) {
	result := make([]rune, 0)
	for _, str := range railTable {
		result = append(result, str...)
	}

	copy(dst, result)
}

func (r *RailCrypto) Decode(dst, src []rune) {
	if r.key < 2 {
		return
	}

	var (
		enc    = src
		encLen = len(enc)
		table  = make([][]rune, r.key)

		periodCount = encLen / r.period()
		ost         = encLen % r.period()
	)

	var (
		index  = 0
		strLen = periodCount
	)

	if ost > 0 {
		strLen = periodCount + 1
	}

	table[index] = enc[:strLen]
	enc = enc[strLen:]

	for index = index + 1; index < r.key-1; index++ {
		strLen := 2 * periodCount

		if index < ost {
			if r.key-1-index <= ost-r.key {
				strLen++
			}

			strLen++
		}

		log.Printf("%s", string(enc[:strLen]))
		table[index] = enc[:strLen]
		enc = enc[strLen:]

	}

	log.Printf("%s", string(enc))
	table[index] = enc

	str := r.generateDecoding(encLen, table)

	copy(dst, str)
}

func (r *RailCrypto) period() int {
	return 2 * (r.key - 1)
}

type TableIterator struct {
	table  [][]rune
	key    int
	strLen int

	iterIndex int
	indexes   []int
	height    int
	isDec     bool
}

func (i *TableIterator) inc() {
	i.indexes[i.height]++

	if i.height == 0 {
		i.isDec = false
	}

	if i.height == i.key-1 {
		i.isDec = true
	}

	if i.isDec {
		i.height--
	} else {
		i.height++
	}
}

func (i *TableIterator) char() rune {
	if i.indexes == nil {
		i.indexes = make([]int, i.key)
	}

	ch := i.table[i.height][i.indexes[i.height]]
	i.inc()

	return ch
}

func (i *TableIterator) next() bool {
	if i.iterIndex == i.strLen {
		return false
	}

	i.iterIndex++

	return true
}

func (r *RailCrypto) generateDecoding(strLen int, table [][]rune) []rune {
	var (
		iter = &TableIterator{
			table:  table,
			key:    r.key,
			strLen: strLen,
		}

		result = make([]rune, strLen)
		i      = 0
	)

	for iter.next() {
		result[i] = iter.char()

		i++
	}

	return result
}
