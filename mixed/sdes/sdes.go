package sdes

// SDES proccesing SDES crypto algorithm.
type SDES struct {
	k1, k2 byte

	s1, s2 [][]byte
}

// New create SDES object.
func New(key []byte) SDES {
	block := SDES{}

	block.k1, block.k2 = generateKeys(key)

	block.s1 = [][]byte{
		{1, 0, 3, 2},
		{3, 2, 1, 0},
		{0, 2, 1, 3},
		{3, 1, 3, 2},
	}
	block.s2 = [][]byte{
		{0, 1, 2, 3},
		{2, 0, 1, 3},
		{3, 0, 1, 0},
		{2, 1, 0, 3},
	}

	return block
}

func generateKeys(key []byte) (byte, byte) {
	var (
		key10Rearrange = []int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}
		key8Rearrange  = []int{6, 3, 7, 4, 8, 5, 10, 9}
	)

	// from generate key method
	firstStepResult := rearrange(key10Rearrange, []byte(key), 10)
	left, right := halve10Bit(firstStepResult)
	left, right = rotateLeft5(left, 1), rotateLeft5(right, 1)

	k1 := rearrange(key8Rearrange, contact5Bit(left, right), 10)[0]

	left, right = rotateLeft5(left, 2), rotateLeft5(right, 2)

	k2 := rearrange(key8Rearrange, contact5Bit(left, right), 10)[0]

	return k1, k2
}

// BlockSize returns size of the block data.
func (s SDES) BlockSize() int {
	return 8
}

// Encrypt encrypt src data for the SDES algorithm and write encryption to the dst.
func (s SDES) Encrypt(dst, src []byte) {
	for i, b := range src {
		dst[i] = s.encrypt(b)
	}
}

func (s SDES) encrypt(b byte) byte {
	return s.crypt(b, s.k1, s.k2)
}

func (s SDES) round(b byte, key byte) byte {
	left, right := halve8Bit(b)

	ep := rearrange([]int{4, 1, 2, 3, 2, 3, 4, 1}, []byte{right}, 4)[0]
	xorKey := ep ^ key

	sLeft, sRight := halve8Bit(xorKey)
	slRes, srRes := blockproc(s.s1, sLeft), blockproc(s.s2, sRight)
	sRes := slRes<<2 | srRes

	sRes = rearrange([]int{2, 4, 3, 1}, []byte{sRes}, 4)[0]

	res := sRes ^ left

	return contact4Bit(res, right)
}

func (s SDES) crypt(b byte, k1 byte, k2 byte) byte {
	var (
		startRearrange = []int{2, 6, 3, 1, 4, 8, 5, 7}
		endRearrange   = []int{4, 1, 3, 5, 7, 2, 8, 6}
	)

	start := rearrange(startRearrange, []byte{b}, 8)[0]

	round := s.round(start, k1)

	lround, rround := halve8Bit(round)

	round = s.round(contact4Bit(rround, lround), k2)

	return rearrange(endRearrange, []byte{round}, 8)[0]

}

func blockproc(sblock [][]byte, b byte) byte {
	row, col := getSBlockCord(b)

	return sblock[row][col]
}

func getSBlockCord(b byte) (byte, byte) {
	return 0x3 & (replaceBit(getBit(3, b), 2, 4) | (0x1 & b)), (0x6 & b) >> 1
}

// Decrypt decrypts src encryption to the source text for the SDES algorithm.
func (s SDES) Decrypt(dst, src []byte) {

}

func (s SDES) decrypt(b byte) byte {
	return s.crypt(b, s.k2, s.k1)
}
