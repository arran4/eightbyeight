package main

import (
	"github.com/arran4/eightbyeight"
	"image/color"
	"log"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)

	err := eightbyeight.NewGridBuilder().
		WithTitle("Grid draw test - Red and Blue").
		WithDimensions(16*16/4, 4).
		WithColors([]color.Color{color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}}).
		Save("out.png")

	if err != nil {
		log.Panic(err)
	}
}
