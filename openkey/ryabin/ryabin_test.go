package ryabin

import (
	"crypto/rand"
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGCD(t *testing.T) {
	tt := []struct {
		a, b *big.Int

		got, x, y *big.Int
	}{
		{
			a: new(big.Int).SetInt64(240),
			b: new(big.Int).SetInt64(46),

			got: new(big.Int).SetInt64(2),
			x:   new(big.Int).SetInt64(-9),
			y:   new(big.Int).SetInt64(47),
		},
		{
			a: new(big.Int).SetInt64(11),
			b: new(big.Int).SetInt64(19),

			got: new(big.Int).SetInt64(1),
			x:   new(big.Int).SetInt64(7),
			y:   new(big.Int).SetInt64(-4),
		},
	}

	for _, tc := range tt {
		d, x, y := gcd(tc.a, tc.b)

		if d.Cmp(tc.got) != 0 {
			t.Errorf("d is not equal: actual: %d; excepted: %d;", d.Int64(), tc.got.Int64())
		}

		if x.Cmp(tc.x) != 0 {
			t.Errorf("x is not equal: actual: %d; excepted: %d;", x.Int64(), tc.x.Int64())
		}

		if y.Cmp(tc.y) != 0 {
			t.Errorf("y is not equal: actual: %d; excepted: %d;", y.Int64(), tc.y.Int64())
		}
	}
}

func TestMPow(t *testing.T) {
	tt := []struct {
		p *big.Int
		q *big.Int

		mp, mq *big.Int
	}{
		{
			p: new(big.Int).SetInt64(11),
			q: new(big.Int).SetInt64(19),

			mp: new(big.Int).SetInt64(3),
			mq: new(big.Int).SetInt64(5),
		},
	}

	for _, tc := range tt {
		key := &PrivateKey{
			P: tc.p,
			Q: tc.q,
		}

		mp, mq := key.dPows()

		if tc.mp.Cmp(mp) != 0 {
			t.Fatalf("incorrect mp pow; actual: %v; excepted: %v", mp, tc.mp)
		}

		if tc.mq.Cmp(mq) != 0 {
			t.Fatalf("incorrect mq pow; actual: %v; excepted: %v", mq, tc.mq)
		}
	}
}

func TestM(t *testing.T) {
	tt := []struct {
		cipher, pow, mod *big.Int

		m *big.Int
	}{
		{
			cipher: big.NewInt(80),
			pow:    big.NewInt(3),
			mod:    big.NewInt(11),

			m: big.NewInt(5),
		},
	}

	for _, tc := range tt {
		key := &PrivateKey{}

		m := key.m(tc.cipher, tc.pow, tc.mod)

		if m.Cmp(tc.m) != 0 {
			t.Errorf("m isn't equal; actual: %v, excepted: %v", m, tc.m)
		}
	}
}

func TestDecode(t *testing.T) {
	key, err := GenerateKey(rand.Reader, 256)
	if err != nil {
		t.Fatalf("failed to generate key: %v", err)
	}

	tt := []struct {
		b byte

		key *PrivateKey
	}{
		{
			b:   'a',
			key: key,
		},
	}

	for _, tc := range tt {
		enc := tc.key.PublicKey().encrypt(tc.b)

		got := tc.key.decrypt(enc)

		assert.Equal(t, got, tc.b, "")
	}
}
