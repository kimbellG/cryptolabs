package sdes

import (
	"crypto/cipher"
)

// SDES proccesing SDES crypto algorithm.
type SDES struct {
	k1, k2 byte
}

// New create SDES object.
func New(key []byte) cipher.Block {
	return nil
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
func BlockSize() int {
	return 8
}

// Encrypt encrypt src data for the SDES algorithm and write encryption to the dst.
func (s SDES) Encrypt(dst, src []byte) {

}
