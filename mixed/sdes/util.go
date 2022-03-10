package sdes

func rearrange(order []int, sequence []byte, seqLen int) []byte {
	const (
		byteLen = 8
	)

	resutlLen := len(order)

	resSliceLen := resutlLen / byteLen
	if resutlLen%byteLen > 0 {
		resSliceLen++
	}

	result := make([]byte, resSliceLen)

	for dst, src := range order {
		src, dst = seqLen-src, resutlLen-dst-1
		srcI, dstI := src/byteLen, dst/byteLen

		srcbit := getBit(src%byteLen, sequence[srcI])
		result[dstI] |= replaceBit(srcbit, dst%byteLen, src%byteLen)
	}

	return result
}

func getBit(index int, b byte) byte {
	var result byte = 1 << index

	return b & result
}

func replaceBit(b byte, dst, src int) byte {
	if dst < src {
		return b >> (src - dst)
	}

	return b << (dst - src)
}

func halve10Bit(src []byte) (byte, byte) {
	if len(src) != 2 {
		panic("halve 10 bit: invalid src slice")
	}

	return (0x18 & src[1] << 3) | (src[0]&0xE0)>>5, src[0] & 0x1F
}

func contact5Bit(left, right byte) []byte {
	return []byte{left<<5 | right, 0x03 & left >> 5}
}

func halve8Bit(b byte) (byte, byte) {
	return (0xF0 & b) >> 4, 0xF & b
}

func contact4Bit(left, right byte) byte {
	return left<<4 | right
}

func rotateLeft5(b byte, k int) byte {
	const n = 5

	s := mod(k, n)
	return 0x1F & (b<<s | b>>(n-s))
}

func mod(l int, r int) int {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}

		return x
	}

	s := abs(l) % r
	switch {
	case s == 0:
		return 0
	case l < 0:
		return r - s
	}

	return s
}
