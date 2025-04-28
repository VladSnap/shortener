package helpers

import "testing"

const randStringRues int = 8

func BenchmarkRandStringRunes(b *testing.B) {
	for b.Loop() {
		_, _ = RandStringRunes(randStringRues)
	}
}
