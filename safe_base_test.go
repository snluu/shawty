package main

import (
	"math/rand"
	"testing"
)

// TestSanity makes sure g(f(x)) = x
func TestSanity(t *testing.T) {
	for i := 0; i < 10000; i++ {
		var n = uint64(rand.Int63())
		var safeBase = toSafeBase(n)
		var dec = toDec(safeBase)
		if n != dec {
			t.Errorf("Not sane: n = %d, sb = %s, dec = %d", n, safeBase, dec)
			t.FailNow()
		}
	}
}
