package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"log"
	"os"

	"github.com/arran4/eightbyeight"
)

// DrawFloppy draws a stylized 90s aesthetic Windows 3.11 floppy disk to the given destination.
// It uses multiple ColourSource patterns to simulate plastic, metal, and label textures.
func DrawFloppy(dst *image.RGBA) {
	// Base floppy shape
	// Win 3.11 floppy disks had colored plastic, e.g., bright blue or magenta.
	// Let's use a nice patterned cyan/magenta for the plastic shell.
	plasticPattern := eightbyeight.NewColourSource(110, color.RGBA{0, 255, 255, 255}, color.RGBA{255, 0, 255, 255})

	// Create plastic shape: basic square, minus the top-right notch
	plasticRect := image.Rect(32, 32, 224, 224)
	draw.Draw(dst, plasticRect, plasticPattern, image.Point{}, draw.Over)

	// Top right notch (transparent/background)
	draw.Draw(dst, image.Rect(192, 32, 224, 64), image.NewUniform(color.RGBA{255, 255, 255, 255}), image.Point{}, draw.Src)

	// Metal slider: top middle (shiny metallic pattern)
	metalPattern := eightbyeight.NewColourSource(42, color.RGBA{220, 220, 220, 255}, color.RGBA{160, 160, 160, 255})
	sliderRect := image.Rect(80, 32, 176, 96)
	draw.Draw(dst, sliderRect, metalPattern, image.Point{}, draw.Over)

	// Slider groove/hole
	draw.Draw(dst, image.Rect(144, 48, 160, 80), image.NewUniform(color.Black), image.Point{}, draw.Src)

	// White label: bottom middle
	labelRect := image.Rect(56, 112, 200, 216)
	draw.Draw(dst, labelRect, image.NewUniform(color.White), image.Point{}, draw.Src)

	// Label lines (text simulation)
	draw.Draw(dst, image.Rect(72, 136, 184, 144), image.NewUniform(color.Black), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(72, 160, 184, 168), image.NewUniform(color.Black), image.Point{}, draw.Src)
	draw.Draw(dst, image.Rect(72, 184, 152, 192), image.NewUniform(color.Black), image.Point{}, draw.Src)
}

func ExampleFloppy() {
	dst := image.NewRGBA(image.Rect(0, 0, 256, 256))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	DrawFloppy(dst)

	f, err := os.Create("out_floppy.png")
	if err != nil {
		log.Panic(fmt.Errorf("failed to create out_floppy.png: %w", err))
	}
	defer f.Close()
	if err := png.Encode(f, dst); err != nil {
		log.Panic(fmt.Errorf("failed to encode out_floppy.png: %w", err))
	}
	log.Println("Output saved to out_floppy.png")
}
