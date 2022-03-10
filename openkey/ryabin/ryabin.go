package ryabin

import (
	"crypto/rand"
	"errors"
	"fmt"
	"io"
	"math/big"
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

func split(src []byte, length int, f func(r []byte)) {
	i := 0
	for ; i < len(src)/length; i++ {
		f(src[i*length : (i+1)*length])
	}

	end := src[i*length:]
	f(end)
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

func (k *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	if err := checkPub(k); err != nil {
		return nil, fmt.Errorf("validate public key: %w", err)
	}

	res := make([]byte, 0, 2*len(msg))

	for _, b := range msg {
		res = append(res, k.encode(b)...)
	}

	return res, nil
}

func (k *PublicKey) encode(ch byte) []byte {
	res := new(big.Int).SetBytes([]byte{'*', ch})

	return modpow(res, k.N, 2).Bytes()
}

func modpow(n, mod *big.Int, p int) *big.Int {
	res := new(big.Int).SetInt64(1)

	for i := 0; i < p; i++ {
		res = res.Mul(n, res)
		res = res.Mod(res, mod)
	}

	return res
}

func checkPub(key *PublicKey) error {
	zero := new(big.Int).SetInt64(0)
	if key.N.Cmp(zero) == 0 {
		return errors.New("zero key")
	}

	return nil
}

func (k *PrivateKey) Decrypt(cipher []byte) ([]byte, error) {
	if err := k.checkCipher(cipher); err != nil {
		return nil, fmt.Errorf("check cipher: %w", err)
	}

	msg := make([]byte, 0, len(cipher)/2)

	split(cipher, 2, func(r []byte) {
	})

	return msg, nil
}

// block is a 2 character.
func (k *PrivateKey) decrypt(block []byte) byte {
	cipher := new(big.Int).SetBytes(block)

	mpPow, mqPow := k.dPows()

	mp := new(big.Int).Exp(cipher, mpPow, k.p)
	mq := new(big.Int).Exp(cipher, mqPow, k.q)

}

func (k *PrivateKey) dPows() (mp, mq *big.Int) {
	var (
		one  = new(big.Int).SetInt64(1)
		four = new(big.Int).SetInt64(4)
	)

	mq = new(big.Int).Add(k.q, one)
	mp = new(big.Int).Add(k.p, one)

	return mp.Div(mp, four), mq.Div(mq, four)
}

func (k *PrivateKey) checkCipher(cipher []byte) error {
	if len(cipher)%2 != 0 {
		return fmt.Errorf("cipher length should be 2x")
	}

	return nil
}
