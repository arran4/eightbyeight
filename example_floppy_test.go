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
	dst := image.NewRGBA(image.Rect(0, 0, 256, 256))

	// Fill background with white
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	// Draw plastic base
	plasticPattern := eightbyeight.NewColourSource(110, color.RGBA{0, 255, 255, 255}, color.RGBA{255, 0, 255, 255})
	draw.Draw(dst, image.Rect(32, 32, 224, 224), plasticPattern, image.Point{}, draw.Over)
	draw.Draw(dst, image.Rect(192, 32, 224, 64), image.NewUniform(color.White), image.Point{}, draw.Src) // Notch

	// Metal slider
	metalPattern := eightbyeight.NewColourSource(42, color.RGBA{220, 220, 220, 255}, color.RGBA{160, 160, 160, 255})
	draw.Draw(dst, image.Rect(80, 32, 176, 96), metalPattern, image.Point{}, draw.Over)
	draw.Draw(dst, image.Rect(144, 48, 160, 80), image.NewUniform(color.Black), image.Point{}, draw.Src) // Hole

	// Label
	draw.Draw(dst, image.Rect(56, 112, 200, 216), image.NewUniform(color.White), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(72, 136, 184, 144), image.NewUniform(color.Black), image.Point{}, draw.Src) // Text

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
