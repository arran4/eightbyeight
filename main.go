package main

import (
	"fmt"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"image/draw"
	"log"
	"os"
)

func IntMax(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

func Int26_6Max(i1, i2 fixed.Int26_6) fixed.Int26_6 {
	if i1 > i2 {
		return i1
	}
	return i2
}

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	log.Printf("Setup")
	pal := []color.Color{
		color.White,
		color.Black,
	}
	fc, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Panicf("Font parse error: %#v", err)
	}
	fontFace := truetype.NewFace(fc, &truetype.Options{
		Size: 16,
		DPI:  150,
	})
	fontHeight := fontFace.Metrics().Ascent
	lineHeight := fontFace.Metrics().Height + fontFace.Metrics().Descent

	labelBounds, _ := font.BoundString(fontFace, fmt.Sprintf("__%d__", 255))

	totalSize := image.Rect(0, 0, IntMax(labelBounds.Max.X.Ceil(), 64)*16, (lineHeight.Ceil()+64)*17+lineHeight.Ceil())
	i := image.NewPaletted(totalSize, pal)
	log.Print("Adding header")
	d := &font.Drawer{
		Dst:  i,
		Src:  image.NewUniform(color.Black),
		Face: fontFace,
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fontHeight,
		},
	}
	d.DrawString("Grid draw test - black and white")
	log.Printf("Drawing grid with labels")
	for y := 0; y < 16; y++ {
		yTop := lineHeight.Ceil() + (lineHeight.Ceil()+64)*(y)
		for x := 0; x < 16; x++ {
			xLeft := IntMax(labelBounds.Max.X.Ceil(), 64) * x
			r := image.Rect(
				xLeft,
				yTop,
				IntMax(labelBounds.Max.X.Ceil(), 64)*(x+1)-1,
				yTop+64-1,
			)
			d.Dot.X = fixed.I(IntMax(labelBounds.Max.X.Ceil(), 64) * (x))
			d.Dot.Y = fixed.I(yTop + 64 - 1 + fontHeight.Ceil())
			dy, dx := r.Dy(), r.Dx()
			if dy > 64 {
				r.Min.Y += (dy - 64) / 2
				r.Max.Y -= (dy - 64) / 2
			}
			if dx > 64 {
				r.Min.X += (dx - 64) / 2
				r.Max.X -= (dx - 64) / 2
			}
			d.DrawString(fmt.Sprintf("  %d", x+y*16))
			if true {
				draw.Draw(i, r, image.NewUniform(color.Black), image.Point{}, draw.Src)
			}
		}
	}
	log.Print("Writing file")
	f, err := os.Create("out.bmp")
	if err != nil {
		log.Panicf("Failed to create file: %#v", err)
	}
	if err := bmp.Encode(f, i); err != nil {
		log.Panic(err)
	}
	log.Printf("Done")
}
