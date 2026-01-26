package main

import (
	"github.com/arran4/eightbyeight"
	"image/color"
	"log"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	err := eightbyeight.NewGridBuilder().
		WithTitle("Grid draw test - black and white").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{color.White, color.Black}).
		Save("out.bmp")

	if err != nil {
		log.Panic(err)
	}
}
