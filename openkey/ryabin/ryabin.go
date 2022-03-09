package ryabin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
	"unicode/utf8"
)

// PublicKey is a public key for the ryabin algo.
type PublicKey struct {
	N *big.Int
}

// PrivateKey is a private key for the ryabin algo.
type PrivateKey struct {
	public *PublicKey
	p, q   *big.Int
}

// GenerateKey creates private key for ryabin algo.
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
	var (
		key = &PrivateKey{}
		err error
	)

	key.p, key.q, err = genPQ(random, bits)
	if err != nil {
		return nil, fmt.Errorf("generate pq key: %w", err)
	}

	key.public = &PublicKey{
		N: new(big.Int).Mul(key.p, key.q),
	}

	return key, nil
}

func genPQ(random io.Reader, bits int) (p *big.Int, q *big.Int, err error) {
	module := big.NewInt(4)

	p, err = rand.Prime(random, bits)
	if err != nil {
		return nil, nil, fmt.Errorf("generate p key: %w", err)
	}

	for {
		q, err = rand.Prime(random, bits)
		if err != nil {
			return nil, nil, fmt.Errorf("generate q key: %w", err)
		}

		if new(big.Int).Mod(p, module).Int64() == 3 &&
			new(big.Int).Mod(q, module).Int64() == 3 {
			return p, q, nil
		}
	}
}

// PublicKey returns public key for the private key.
func (k *PrivateKey) PublicKey() *PublicKey {
	return k.public
}

func (k *PublicKey) Encrypt(src []rune) ([]byte, error) {

	if err := checkPub(k); err != nil {
		return nil, fmt.Errorf("validate public key: %w", err)
	}

	for _, r := range src {
		cipherByte := make([]byte, 4)

		n := utf8.EncodeRune(cipherByte, r)

		enc := k.encode(cipherByte[:n])
	}
}

func (k *PublicKey) encode(cipher []byte) []byte {
	res := new(big.Int).SetBytes(cipher)

}

func checkPub(key *PublicKey) error {
	zero := new(big.Int).SetInt64(0)
	if key.N.Cmp(zero) == 0 {
		return errors.New("zero key")
	}

	return nil
}
