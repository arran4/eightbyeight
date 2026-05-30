package eightbyeight_test

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/arran4/eightbyeight"
)

func ExampleColourSource_floppy() {
	// Create a new RGBA image to draw into
	dst := image.NewRGBA(image.Rect(0, 0, 128, 128))

	// Fill background with white
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	// Draw plastic base
	plasticPattern := eightbyeight.NewColourSource(110, color.RGBA{0, 255, 255, 255}, color.RGBA{255, 0, 255, 255})
	draw.Draw(dst, image.Rect(16, 16, 112, 112), plasticPattern, image.Point{}, draw.Over)
	draw.Draw(dst, image.Rect(96, 16, 112, 32), image.NewUniform(color.White), image.Point{}, draw.Src) // Notch

	// Metal slider
	metalPattern := eightbyeight.NewColourSource(42, color.RGBA{220, 220, 220, 255}, color.RGBA{160, 160, 160, 255})
	draw.Draw(dst, image.Rect(40, 16, 88, 48), metalPattern, image.Point{}, draw.Over)
	draw.Draw(dst, image.Rect(72, 24, 80, 40), image.NewUniform(color.Black), image.Point{}, draw.Src) // Hole

	// Label
	draw.Draw(dst, image.Rect(28, 56, 100, 108), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(36, 68, 92, 72), image.NewUniform(color.Black), image.Point{}, draw.Src) // Text

	f, err := os.Create("out_floppy.png")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer f.Close()

	if err := png.Encode(f, dst); err != nil {
		fmt.Println("Error:", err)
		return
	}

	fmt.Println("Output saved to out_floppy.png")
	// Output: Output saved to out_floppy.png
}
