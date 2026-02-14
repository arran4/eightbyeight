package eightbyeight

import (
	"image"
	"image/color"
	"math"
	"testing"
)

// LegacyColourSource implementation for regression testing
type LegacyColourSource struct {
	colors []color.Color
	mode   int
	sz     int
}

func NewLegacyColourSource(mode int, colors ...color.Color) *LegacyColourSource {
	sz := 8
	return &LegacyColourSource{
		mode:   mode,
		colors: colors,
		sz:     sz,
	}
}

func (cs *LegacyColourSource) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{-1e9, -1e9},
		image.Point{1e9, 1e9},
	}
}

func (cs *LegacyColourSource) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color { return c })
}

func (cs *LegacyColourSource) At(x, y int) color.Color {
	south := []int{
		1,
		0,
		0,
		0,
	}
	n := cs.mode
	for i := 0; i < 4 && n > 0; i++ {
		south[i] = n % 4
		n /= 4
	}

	xp := x % cs.sz
	dp := (y + cs.sz - xp) % cs.sz

	if cs.mode >= xp {
		xp = 3 - (3 - xp)
		if xp < 0 {
			xp = -xp
		}
		sv := south[xp%len(south)]
		sv = int(math.Pow(float64(2), float64(4-sv)))
		if sv > 0 && dp%sv == 0 {
			// Use the second color (index 1) for "foreground"/black if available
			if len(cs.colors) > 1 {
				return cs.colors[1]
			}
			return color.Black
		}
	}

	// Use the first color (index 0) for "background"/white if available
	if len(cs.colors) > 0 {
		return cs.colors[0]
	}
	return color.White
}

func TestColourSource_Correctness(t *testing.T) {
	// Test range of modes and coordinates
	// Modes can be quite large, but let's test a reasonable range that covers the logic
	// The mode affects the 'south' array calculation.
	// Since south is 4 integers derived from mode, let's test enough modes to vary south.

	colors := []color.Color{color.White, color.Black}

	for mode := 0; mode < 512; mode++ {
		optimized := NewColourSource(mode, colors...)
		legacy := NewLegacyColourSource(mode, colors...)

		// Test a grid larger than sz (8) to ensure tiling works
		for y := 0; y < 16; y++ {
			for x := 0; x < 16; x++ {
				optColor := optimized.At(x, y)
				legColor := legacy.At(x, y)

				r1, g1, b1, a1 := optColor.RGBA()
				r2, g2, b2, a2 := legColor.RGBA()

				if r1 != r2 || g1 != g2 || b1 != b2 || a1 != a2 {
					t.Errorf("Mismatch at mode=%d, x=%d, y=%d. Optimized: %v, Legacy: %v",
						mode, x, y, optColor, legColor)
				}
			}
		}
	}
}
