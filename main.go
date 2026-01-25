package main

import (
	"image/color"
	"log"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	err := NewGridBuilder().
		WithTitle("Grid draw test - black and white").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{color.White, color.Black}).
		Save("out.bmp")

	if err != nil {
		log.Panic(err)
	}
}
