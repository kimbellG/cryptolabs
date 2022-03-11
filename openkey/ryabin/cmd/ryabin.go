package main

import (
	"crypto/rand"
	"encoding/gob"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/kimbellG/cryptolabs/openkey/ryabin"
)

var (
	secret = flag.String("s", "", "secret key filepath")

	keygen = flag.Bool("k", false, "generate key")
	bits   = flag.Int("b", 0, "key bits count")

	decrypt = flag.Bool("d", false, "decrypt mode")
	public  = flag.String("p", "", "public key filepath")
)

func main() {
	flag.Parse()

	m()
}

func m() {
	switch {
	case *keygen:
		if *bits == 0 {
			log.Fatalf("ryabin: key size is not declare")
		}

		keygenerate(*bits)

		return

	case *decrypt:
		if *secret == "" {
			log.Fatalf("ryabin: secret file is not declare")
		}

		for _, arg := range flag.Args() {
			if err := decryptFile(arg, arg+".decrypt", *secret); err != nil {
				log.Fatalf("ryabin: decrypt file: %v", err)
			}
		}

		return
	default:
		if *public == "" {
			log.Fatalf("ryabin: public key file is not declare")
		}

		for _, arg := range flag.Args() {
			if err := encryptFile(arg, arg+".encrypt", *public); err != nil {
				log.Fatalf("ryabin: encrypt file: %v", err)
			}
		}
	}
}

func keygenerate(bits int) {
	key, err := ryabin.GenerateKey(rand.Reader, bits)
	if err != nil {
		log.Fatalf("ryabin: generate key: %v", err)
	}

	log.Printf("p: %v, q: %v, n: %v", key.P, key.Q, key.Public.N)

	pub := key.PublicKey()

	if err := saveGobKey("private.key", key); err != nil {
		log.Fatalf("ryabin: save private key: %v", err)
	}

	if err := saveGobKey("public.key", pub); err != nil {
		log.Fatalf("ryabin: save public key: %v", err)
	}
}

func encryptFile(filename, dest string, publicfile string) error {
	pub := &ryabin.PublicKey{}

	err := getGobKey(publicfile, pub)
	if err != nil {
		return fmt.Errorf("get public key from file: %w", err)
	}

	log.Printf("n: %v", pub.N)

	in, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open decrypt filename: %w", err)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}

	buf := make([]byte, 200)

	var n int
	for n, err = in.Read(buf); err == nil; n, err = in.Read(buf) {
		dst, err := pub.Encrypt(buf[:n])
		if err != nil {
			return fmt.Errorf("encrypt buffer: %w", err)
		}

		_, werr := out.Write(dst)
		if werr != nil {
			log.Printf("%s: could not write to the output file: %v", flag.Arg(0), werr)
		}
	}

	if !errors.Is(err, io.EOF) {
		return fmt.Errorf("read from cipher file: %w", err)
	}

	if n == 0 {
		return nil
	}

	dst, err := pub.Encrypt(buf[:n])
	if err != nil {
		return fmt.Errorf("encrypt last block: %w", err)
	}

	_, err = out.Write(dst)
	if err != nil {
		return fmt.Errorf("write last block to the file: %w", err)
	}

	return nil
}

func decryptFile(filename, dest string, secretfile string) error {
	key := &ryabin.PrivateKey{}

	err := getGobKey(secretfile, key)
	if err != nil {
		return fmt.Errorf("could not get key from file: %w", err)
	}

	log.Printf("p: %v, q: %v, n: %v", key.P, key.Q, key.Public.N)

	in, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("open decrypt filename: %w", err)
	}

	out, err := os.Create(dest)
	if err != nil {
		return fmt.Errorf("create destination file: %w", err)
	}

	buf := make([]byte, key.PublicKey().BlockSize())

	for _, err := in.Read(buf); !errors.Is(err, io.EOF); _, err = in.Read(buf) {
		dst, err := key.Decrypt(buf)
		if err != nil {
			return fmt.Errorf("read cipher file: %w", err)
		}

		_, werr := out.Write(dst)
		if werr != nil {
			log.Printf("%s: could not write to the output file: %v", flag.Arg(0), werr)
		}
	}

	return nil
}

func saveGobKey(filename string, key interface{}) error {
	out, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("could not create file for the key: %w", err)
	}
	defer out.Close()

	encoder := gob.NewEncoder(out)

	if err := encoder.Encode(key); err != nil {
		return fmt.Errorf("could not encode key: %w", err)
	}

	return nil
}

func getGobKey(filename string, key interface{}) error {
	in, err := os.Open(filename)
	if err != nil {
		return fmt.Errorf("could not open key file: %w", err)
	}
	defer in.Close()

	decoder := gob.NewDecoder(in)

	if err := decoder.Decode(key); err != nil {
		return fmt.Errorf("could not decode key from file: %w", err)
	}

	return nil
}
