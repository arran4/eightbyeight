package main

import (
	"image/color"
	"testing"
)

func TestGridBuilder_Generate(t *testing.T) {
	builder := NewGridBuilder().
		WithTitle("Test Grid").
		WithDimensions(2, 2).
		WithColors([]color.Color{color.White, color.Black})

	img := builder.Generate()

	if img == nil {
		t.Fatal("Generate() returned nil")
	}

	bounds := img.Bounds()
	if bounds.Empty() {
		t.Error("Generated image has empty bounds")
	}

	// We can check if dimensions are roughly what we expect.
	// Rows=2, Cols=2.
	// Logic in builder:
	// lines = 2, lineLength = 2.
	// totalSize calculated.
	// Just ensuring it runs without panic and returns a valid image is good for now.
}

func TestNewColourSource(t *testing.T) {
	cs := NewColourSource(0, color.White, color.Black)
	if cs == nil {
		t.Fatal("NewColourSource returned nil")
	}
	// Check At(0,0) returns a color
	c := cs.At(0, 0)
	if c == nil {
		t.Error("At(0,0) returned nil color")
	}
}
