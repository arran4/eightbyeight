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

func ExampleColourSource_floppy() {
	// Create a new RGBA image to draw into
	dst := image.NewRGBA(image.Rect(0, 0, 16, 16))

	// Fill background with white
	draw.Draw(dst, dst.Bounds(), image.NewUniform(color.White), image.Point{}, draw.Src)

	// Create a ColourSource pattern
	pattern := eightbyeight.NewColourSource(110, color.White, color.Black)

	// Define the mask for the floppy disk
	mask := &FloppyMask{}

	// Fill the destination image using the mask and the pattern
	draw.DrawMask(dst, dst.Bounds(), pattern, image.Point{}, mask, image.Point{}, draw.Over)

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
