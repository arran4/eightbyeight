package eightbyeight

import (
	"image/color"
	"io/ioutil"
	"log"
	"testing"
)

func BenchmarkGridBuilder_Generate(b *testing.B) {
	// Suppress logging
	log.SetOutput(ioutil.Discard)

	// Setup a builder with a larger grid to emphasize the loop performance
	builder := NewGridBuilder().
		WithTitle("Benchmark Grid").
		WithDimensions(50, 20). // 1000 cells
		WithColors([]color.Color{color.White, color.Black})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		builder.Generate()
	}
}
