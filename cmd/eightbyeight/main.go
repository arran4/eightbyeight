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

	// Example 5: Stylized Floppy Disk
	// Demonstrates drawing a stylized 90s aesthetic floppy using multiple image/draw operations and patterns
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
