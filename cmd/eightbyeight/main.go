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

// FloppyMask defines a simple 16x16 mask for a floppy disk.
// 1 represents the floppy body, 0 represents transparent areas (like the label and the slider).
var floppyMaskData = []byte{
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,1,1,0,0,0,0,0,0,0,0,1,1,0,0,
	0,0,1,1,0,0,0,0,0,0,0,0,1,1,0,0,
	0,0,1,1,0,0,0,0,0,0,0,0,1,1,0,0,
	0,0,1,1,0,0,0,0,0,0,0,0,1,1,0,0,
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,1,1,1,1,0,0,0,0,1,1,1,1,0,0,
	0,0,1,1,1,1,0,0,0,0,1,1,1,1,0,0,
	0,0,1,1,1,1,0,0,0,0,1,1,1,1,0,0,
	0,0,1,1,1,1,0,0,0,0,1,1,1,1,0,0,
	0,0,1,1,1,1,0,0,0,0,1,1,1,1,0,0,
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,1,1,1,1,1,1,1,1,1,1,1,1,0,0,
	0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,0,
}

type FloppyMask struct{}

func (m *FloppyMask) ColorModel() color.Model { return color.AlphaModel }
func (m *FloppyMask) Bounds() image.Rectangle { return image.Rect(0, 0, 16, 16) }
func (m *FloppyMask) At(x, y int) color.Color {
	if x < 0 || x >= 16 || y < 0 || y >= 16 {
		return color.Alpha{0}
	}
	if floppyMaskData[y*16+x] == 1 {
		return color.Alpha{255}
	}
	return color.Alpha{0}
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

	// Example 5: Floppy Disk with draw.DrawMask
	// Demonstrates using a mask to draw a patterned shape
	dst := image.NewRGBA(image.Rect(0, 0, 16, 16))
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)
	pattern := eightbyeight.NewColourSource(110, color.White, color.Black)
	draw.DrawMask(dst, dst.Bounds(), pattern, image.Point{}, &FloppyMask{}, image.Point{}, draw.Over)

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
