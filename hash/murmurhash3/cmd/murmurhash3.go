package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kimbellG/cryptolabs/hash/murmurhash3"
)

func main() {
	const seed = 21

	var (
		buf = make([]byte, 4*1024)

		hash = murmurhash3.NewHash32(21)

		n   int
		err error
	)

	for n, err = os.Stdin.Read(buf); err == nil; n, err = os.Stdin.Read(buf[4:]) {
		_, err := hash.Write(buf[:n])
		if err != nil {
			log.Fatalf("%s: hash iteration: %v", os.Args[0], err)
		}
	}

	if !errors.Is(err, io.EOF) {
		log.Printf("%s: failed to read from stdin: %v", os.Args[0], err)
	}

	if n > 0 {
		_, err := hash.Write(buf[:n])
		if err != nil {
			log.Fatalf("%s: last hash iteration: %v", os.Args[0], err)
		}
	}

	os.Stdout.WriteString(fmt.Sprintf("hash32: 0x%x\n", hash.SumAndClean()))

}
