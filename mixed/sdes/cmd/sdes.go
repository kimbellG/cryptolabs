package main

import (
	"errors"
	"flag"
	"io"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/kimbellG/cryptolabs/mixed/sdes"
)

var (
	key1   = flag.Uint("k1", uint(generateKey()[1]), "first byte of the key")
	key2   = flag.Uint("k2", uint(generateKey()[0]), "second byte of the key")
	output = flag.String("o", "a.out", "output file")
	isDecr = flag.Bool("e", false, "encrypt")
)

func main() {
	flag.Parse()

	for _, arg := range flag.Args() {
		file, err := os.Open(arg)
		if err != nil {
			log.Fatalf("%s: could not open file %s: %v", flag.Arg(0), arg, err)
		}
		defer file.Close()

		out, err := os.Create(*output)
		if err != nil {
			log.Fatalf("%s: could not create output file(%s): %v", flag.Arg(0), *output, err)
		}
		defer out.Close()

		key := []byte{byte(*key1), byte(*key2)}

		log.Printf("for file %s: key: %b %b. Output: %s", arg, key[1], key[0], *output)

		var (
			buf = make([]byte, 200)
			des = sdes.New(key)
		)

		for _, err := file.Read(buf); !errors.Is(err, io.EOF); _, err = file.Read(buf) {
			dst := make([]byte, 200)

			if *isDecr {
				des.Decrypt(dst, buf)
			} else {
				des.Encrypt(dst, buf)
			}

			_, werr := out.Write(dst)
			if werr != nil {
				log.Printf("%s: could not write to the output file: %v", flag.Arg(0), werr)
			}
		}

	}

}

func generateKey() []byte {
	rand.Seed(time.Now().UnixNano())
	return []byte{byte(rand.Intn(4)), byte(rand.Intn(256))}
}
