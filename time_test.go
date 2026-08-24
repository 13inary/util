package util

import (
	"testing"
	"time"
)

// 16.37 ns/op
func BenchmarkIntTime(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		year, month, day := now.Date()
		now = time.Date(year, month, day, 0, 0, 0, 0, now.Location())
	}
}

// 0.2216 ns/op
func BenchmarkUnixTime(b *testing.B) {
	now := time.Now()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		db := now.Unix()
		now = time.Unix(db, 0)
	}
}
