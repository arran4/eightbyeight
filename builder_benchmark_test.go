package main

import (
	"image/color"
	"testing"
)

func BenchmarkColourSource_At(b *testing.B) {
	cs := NewColourSource(123, color.White, color.Black)
	var c color.Color

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		c = cs.At(i%64, i%64)
	}
	_ = c
}
