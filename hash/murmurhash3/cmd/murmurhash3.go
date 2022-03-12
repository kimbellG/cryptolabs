package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"unsafe"

	"github.com/kimbellG/cryptolabs/hash/murmurhash3"
)

func main() {
	const seed = 21

	var (
		buf = make([]byte, 4*1024)

		hash uint32

		n   int
		err error
	)

	for n, err = os.Stdin.Read(buf); err == nil; n, err = os.Stdin.Read(buf[4:]) {
		var hashBytes []byte

		if hash != 0 {
			hashBytes = (*[4]byte)(unsafe.Pointer(&hash))[:]
			copy(buf[0:4], hashBytes)
		}

		hash = murmurhash3.Sum32WithSeed(buf[:n+len(hashBytes)], seed)
	}

	if !errors.Is(err, io.EOF) {
		log.Printf("%s: failed to read from stdin: %v", os.Args[0], err)
	}

	if n > 0 {

		hashBytes := (*[4]byte)(unsafe.Pointer(&hash))[:]
		copy(buf[0:4], hashBytes)

		hash = murmurhash3.Sum32WithSeed(buf[:n+len(hashBytes)], seed)
	}

	os.Stdout.WriteString(fmt.Sprintf("hash32: 0x%x\n", hash))

}
