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
	Public *PublicKey
	P, Q   *big.Int
}

// GenerateKey creates private key for ryabin algo.
func GenerateKey(random io.Reader, bits int) (*PrivateKey, error) {
	var (
		key = &PrivateKey{}
		err error
	)

	key.P, key.Q, err = genPQ(random, bits)
	if err != nil {
		return nil, fmt.Errorf("generate pq key: %w", err)
	}

	key.Public = &PublicKey{
		N: new(big.Int).Mul(key.P, key.Q),
	}

	return key, nil
}

func split(src []byte, length int, f func(r []byte)) {
	i := 0
	for ; i < len(src)/length; i++ {
		f(src[i*length : (i+1)*length])
	}

	end := src[i*length:]
	if len(end) != 0 {
		f(end)
	}
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

		if p.Cmp(q) == 0 {
			continue
		}

		if new(big.Int).Mod(p, module).Int64() == 3 &&
			new(big.Int).Mod(q, module).Int64() == 3 {
			return p, q, nil
		}
	}
}

// PublicKey returns public key for the private key.
func (k *PrivateKey) PublicKey() *PublicKey {
	return k.Public
}

// Encrypt ....
func (k *PublicKey) Encrypt(msg []byte) ([]byte, error) {
	if err := checkPub(k); err != nil {
		return nil, fmt.Errorf("validate public key: %w", err)
	}

	block := make([]byte, k.BlockSize())

	res := make([]byte, 0, len(block)*len(msg))

	for _, b := range msg {
		r := k.encrypt(b)
		copy(block[len(block)-len(r):], r)
		res = append(res, block...)

		clear(block)
	}

	return res, nil
}

func clear(block []byte) {
	for i := range block {
		block[i] = 0
	}
}

// BlockSize returns size of the encryption block.
func (k *PublicKey) BlockSize() int {
	const byteSize = 8

	size := k.N.BitLen() / 8
	if k.N.BitLen()%byteSize > 0 {
		size++
	}

	return size
}

func (k *PublicKey) encrypt(ch byte) []byte {
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

// Decrypt ....
func (k *PrivateKey) Decrypt(cipher []byte) ([]byte, error) {
	if err := k.checkCipher(cipher); err != nil {
		return nil, fmt.Errorf("check cipher: %w", err)
	}

	blockSize := k.Public.BlockSize()
	msg := make([]byte, 0, len(cipher)/blockSize)

	split(cipher, blockSize, func(r []byte) {
		b := k.decrypt(r)

		msg = append(msg, b)

	})

	return msg, nil
}

// block is a 2 character.
func (k *PrivateKey) decrypt(block []byte) byte {
	cipher := new(big.Int).SetBytes(block)

	result := k.getMsgVariant(cipher)

	for _, res := range result {
		rb := res.Bytes()

		if len(rb) == 0 {
			continue
		}

		if rb[0] == '*' {
			return rb[1]
		}
	}

	panic("block not found")
}

func (k *PrivateKey) getMsgVariant(cipher *big.Int) []*big.Int {
	mpPow, mqPow := k.dPows()

	mp := k.m(cipher, mpPow, k.P)
	mq := k.m(cipher, mqPow, k.Q)

	var yp, yq *big.Int
	if k.P.Cmp(k.Q) == 1 || k.P.Cmp(k.Q) == 0 {
		_, yp, yq = gcd(k.P, k.Q)
	} else {
		_, yq, yp = gcd(k.Q, k.P)
	}

	ypm := new(big.Int).Mul(yp, new(big.Int).Mul(k.P, mq))
	yqm := new(big.Int).Mul(yq, new(big.Int).Mul(k.Q, mp))

	result := make([]*big.Int, 4)
	result[0] = new(big.Int).Add(ypm, yqm)
	result[0].Mod(result[0], k.Public.N)

	result[1] = new(big.Int).Sub(k.Public.N, result[0])

	result[2] = new(big.Int).Sub(ypm, yqm)
	result[2].Mod(result[2], k.Public.N)

	result[3] = new(big.Int).Sub(k.Public.N, result[2])

	return result
}

func (k *PrivateKey) m(cipher, pow, prime *big.Int) *big.Int {
	return new(big.Int).Exp(cipher, pow, prime)
}

func (k *PrivateKey) dPows() (mp, mq *big.Int) {
	var (
		one  = new(big.Int).SetInt64(1)
		four = new(big.Int).SetInt64(4)
	)

	mq = new(big.Int).Add(k.Q, one)
	mp = new(big.Int).Add(k.P, one)

	return mp.Div(mp, four), mq.Div(mq, four)
}

func (k *PrivateKey) checkCipher(cipher []byte) error {
	if len(cipher)%k.Public.BlockSize() != 0 {
		return fmt.Errorf("cipher length should be 2x")
	}

	return nil
}

var zeroBig = new(big.Int).SetInt64(0)

func gcd(a, b *big.Int) (d, x, y *big.Int) {
	if b.Cmp(zeroBig) == 0 {
		return a, new(big.Int).SetInt64(1), new(big.Int).SetInt64(0)
	}

	q, m := new(big.Int).DivMod(a, b, new(big.Int))

	d, x1, y1 := gcd(b, m)

	return d, new(big.Int).Set(y1), new(big.Int).Sub(x1, q.Mul(q, y1))
}
