package helpers

import "testing"

const randStringRues int = 8

func BenchmarkRandStringRunes(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = RandStringRunes(randStringRues)
	}
}
