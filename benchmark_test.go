package main

import (
	"testing"
)

func BenchmarkGenerate(b *testing.B) {
	builder := NewGridBuilder().
		WithTitle("Benchmark Grid").
		WithDimensions(10, 10)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Generate()
	}
}
