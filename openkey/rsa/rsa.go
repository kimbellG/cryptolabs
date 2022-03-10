package rsa

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
)

type PublicKey struct {
	E int
	N *big.Int
}

type PrivateKey struct {
}

func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
	var (
		one = new(big.Int).SetInt64(1)
	)

	if bits <= 2 {
		return nil, errors.New("bit cound less or equal 2")
	}

	var (
		p, q *big.Int
		err  error
	)

	for {
		p, err = rand.Prime(random, bits)
		if err != nil {
			return nil, fmt.Errorf("create prime p: %w", err)
		}

		q, err = rand.Prime(random, bits)
		if err != nil {
			return nil, fmt.Errorf("create prime q: %w", err)
		}

		if p.Cmp(q) != 0 {
			break
		}
	}

	r := new(big.Int).Mul(p, q)
	fi := new(big.Int).Mul(new(big.Int).Sub(p, one), new(big.Int).Sub(q, one))

}
