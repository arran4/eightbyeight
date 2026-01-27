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
}
