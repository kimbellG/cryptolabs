package sdes

import (
	"crypto/cipher"
)

type SDES struct {
	k1, k2 byte
}

func New(key []byte) cipher.Block {
	return nil
}

func generateKeys(key []byte) (byte, byte) {
	//	var (
	//		key10Rearrange = []int{3, 5, 2, 7, 4, 10, 1, 9, 8, 6}
	//		key8Rearrange  = []int{6, 3, 7, 4, 8, 5, 10, 9}
	//	)

	// from generate key method
	//	firstStepResult := rearrange(key10Rearrange, []byte(key))

	return 0, 0
}

func BlockSize() int {
	return 8
}

func (s SDES) Encrypt(dst, src []byte) {

}
