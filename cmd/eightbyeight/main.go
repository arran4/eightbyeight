package main

import (
	"github.com/arran4/eightbyeight"
	"image/color"
	"log"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	// Example 1: Classic Black on White
	// Good for checking basic patterns clearly.
	if err := eightbyeight.NewGridBuilder().
		WithTitle("Classic - Black on White").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{color.White, color.Black}).
		Save("out_bw.png"); err != nil {
		log.Panic(err)
	}

	// Example 2: Terminal Style (Green on Black)
	// A classic "matrix" or terminal look.
	if err := eightbyeight.NewGridBuilder().
		WithTitle("Terminal - Green on Black").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{color.Black, color.RGBA{0, 255, 0, 255}}).
		Save("out_terminal.png"); err != nil {
		log.Panic(err)
	}

	// Example 3: Solarized (Beige and Dark Teal)
	// A lower contrast, pleasing combination.
	solarizedBg := color.RGBA{253, 246, 227, 255}
	solarizedFg := color.RGBA{7, 54, 66, 255}
	if err := eightbyeight.NewGridBuilder().
		WithTitle("Solarized Light").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{solarizedBg, solarizedFg}).
		Save("out_solarized.png"); err != nil {
		log.Panic(err)
	}

	// Example 4: CGA Palette Mixing
	// Demonstrates how dithering can create perceived intermediate colors
	// using a limited 16-color CGA-inspired palette.
	cgaPalette := []color.Color{
		color.RGBA{0x00, 0x00, 0x00, 0xff}, // 0: Black
		color.RGBA{0x00, 0x00, 0xAA, 0xff}, // 1: Blue
		color.RGBA{0x00, 0xAA, 0x00, 0xff}, // 2: Green
		color.RGBA{0x00, 0xAA, 0xAA, 0xff}, // 3: Cyan
		color.RGBA{0xAA, 0x00, 0x00, 0xff}, // 4: Red
		color.RGBA{0xAA, 0x00, 0xAA, 0xff}, // 5: Magenta
		color.RGBA{0xAA, 0x55, 0x00, 0xff}, // 6: Brown
		color.RGBA{0xAA, 0xAA, 0xAA, 0xff}, // 7: Light Gray
		color.RGBA{0x55, 0x55, 0x55, 0xff}, // 8: Dark Gray
		color.RGBA{0x55, 0x55, 0xFF, 0xff}, // 9: Light Blue
		color.RGBA{0x55, 0xFF, 0x55, 0xff}, // 10: Light Green
		color.RGBA{0x55, 0xFF, 0xFF, 0xff}, // 11: Light Cyan
		color.RGBA{0xFF, 0x55, 0x55, 0xff}, // 12: Light Red
		color.RGBA{0xFF, 0x55, 0xFF, 0xff}, // 13: Light Magenta
		color.RGBA{0xFF, 0xFF, 0x55, 0xff}, // 14: Yellow
		color.RGBA{0xFF, 0xFF, 0xFF, 0xff}, // 15: White
	}
	if err := eightbyeight.NewGridBuilder().
		WithTitle("CGA Color Mixing").
		WithDimensions(16, 16).
		WithColors(cgaPalette).
		Save("out_mixing.png"); err != nil {
		log.Panic(err)
	}
}
