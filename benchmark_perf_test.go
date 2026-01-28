package main

import (
	"image/color"
	"testing"
)

func BenchmarkColourSource_At(b *testing.B) {
	cs := NewColourSource(12345, color.White, color.Black)
	// We iterate over a range of x, y to simulate usage, but simplified.
	// Since the calculation inside At depends on cs.mode (constant here) and x, y.
	// The part we are optimizing is the `south` array calculation which is constant per `cs` instance.
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cs.At(i%100, i%100)
	}
}
