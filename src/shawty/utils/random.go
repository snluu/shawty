package utils

import (
	crand "crypto/rand"
	mrand "math/rand"
	"time"
)

type Rand interface {
	Byte() byte
	Uint32() uint32
	Uint64() uint64
}

// NewBestRand creates a new BestRand instance
func NewBestRand() *BestRand {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	r.Seed(time.Now().UnixNano())
	return &BestRand{r}
}

// NewFakeRand creates a new FakeRand instance
func NewFakeRand() *FakeRand {
	return &FakeRand{0}
}

// BestRand generates random numbers using crypto/random
// and falls back to using math/rand if fails.
// It implements the Rand interface
type BestRand struct {
	rand *mrand.Rand
}

// getBytes gets n random bytes
func (br *BestRand) getBytes(n int) []byte {
	b := make([]byte, n)
	x, err := crand.Read(b)
	if err != nil || x != n {
		for i := range b {
			b[i] = byte(br.rand.Int31())
		}
	}
	return b
}

// Byte generates a new random byte
func (br *BestRand) Byte() byte {
	return br.getBytes(1)[0]
}

// Uint32 generates a random unsigned 32 bit integer
func (br *BestRand) Uint32() uint32 {
	b := br.getBytes(4)
	n := uint32(0)
	for _, x := range b {
		n <<= 8
		n |= uint32(x)
	}
	return n
}

// Uint64 generates a random unsigned 64 bit integer
func (br *BestRand) Uint64() uint64 {
	b := br.getBytes(4)
	n := uint64(0)
	for _, x := range b {
		n <<= 8
		n |= uint64(x)
	}
	return n
}

// FakeRand does not generate random numbers.
// Rather, it returns the seed value (truncated if necessary)
type FakeRand struct {
	seed uint64
}

func (fr *FakeRand) Seed(n uint64) {
	fr.seed = n
}

func (fr *FakeRand) Byte() byte {
	return byte(fr.seed)
}

func (fr *FakeRand) Uint32() uint32 {
	return uint32(fr.seed)
}

func (fr *FakeRand) Uint64() uint64 {
	return fr.seed
}
