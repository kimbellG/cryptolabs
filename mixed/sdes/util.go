package sdes

func rearrange(order []int, sequence []byte) []byte {
	const (
		byteLen = 8
	)

	seqLen := len(order)
	result := make([]byte, len(sequence))

	for dst, src := range order {
		src, dst = seqLen-src, seqLen-dst-1
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
	} else {
		return b << (dst - src)
	}
}

func halve10Bit(src []byte) (byte, byte) {
	if len(src) != 2 {
		panic("halve 10 bit: invalid src slice")
	}

	return (0x18 & src[1] << 3) | (src[0]&0xE0)>>5, src[0] & 0x1F
}
