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
	const LINES = 16 * 16 / 4
	const LINELENGTH = 4
	fontFace := truetype.NewFace(fc, &truetype.Options{
		Size: 16,
		DPI:  150,
	})
	fontHeight := fontFace.Metrics().Ascent
	lineHeight := fontFace.Metrics().Height + fontFace.Metrics().Descent

	labelBounds, _ := font.BoundString(fontFace, fmt.Sprintf("__%d__", 255))
	title := "Grid draw test - black and white"
	titleBounds, _ := font.BoundString(fontFace, title)

	totalSize := image.Rect(0, 0, IntMax(titleBounds.Max.X.Ceil(), IntMax(labelBounds.Max.X.Ceil(), 64)*LINELENGTH), (lineHeight.Ceil()+64)*(LINES)+lineHeight.Ceil())
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
	d.DrawString(title)
	log.Printf("Drawing grid with labels")
	for y := 0; y < LINES; y++ {
		yTop := lineHeight.Ceil() + (lineHeight.Ceil()+64)*(y)
		for x := 0; x < LINELENGTH; x++ {
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
			mode := x + y*LINELENGTH
			d.DrawString(fmt.Sprintf("  %d", mode))
			if true {
				draw.Draw(i, r, NewColourSource(mode), image.Point{}, draw.Src)
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

func NewColourSource(mode int, colors ...color.Color) image.Image {
	sz := 8
	//template := make([][]bool, sz)
	//for i := 0; i < sz; i++ {
	//	template[i] = make([]bool, sz)
	//}
	//if mode > 0 {
	//	for i := 0; i < sz*sz; i += (sz * sz) / mode {
	//		template[i/sz][i%sz] = true
	//	}
	//}
	return &ColourSource{
		mode:   (mode / (256 / (sz * sz))) % (sz * sz),
		colors: colors,
		//template: template,
		sz: sz,
	}
}

type ColourSource struct {
	colors []color.Color
	mode   int
	//template    [][]bool
	sz int
}

func (cs *ColourSource) ColorModel() color.Model {
	return cs
}

func (cs *ColourSource) Convert(cl color.Color) color.Color {
	return cl
}

func (cs *ColourSource) Bounds() image.Rectangle {
	return image.Rectangle{
		image.Point{-1e9, -1e9},
		image.Point{1e9, 1e9},
	}
}

func (cs *ColourSource) At(x, y int) color.Color {
	//if cs.template[y % cs.sz][x % cs.sz] {
	//	return color.Black
	//}
	//xv := x % (cs.sz)
	//yv := y % (cs.sz)
	//n := yv * cs.sz + xv
	//if cs.mode > 0 && n % (cs.sz*cs.sz / cs.mode) == 0 {
	//	return color.Black
	//}
	n := cs.mode
	south := 0
	east := 0
	if (0b000001 & n) > 0 {
		south |= 0b001
	}
	if (0b000100 & n) > 0 {
		south |= 0b010
	}
	if (0b010000 & n) > 0 {
		south |= 0b100
	}
	if (0b000010 & n) > 0 {
		east |= 0b001
	}
	if (0b001000 & n) > 0 {
		east |= 0b010
	}
	if (0b100000 & n) > 0 {
		east |= 0b100
	}
	//southEast := 0
	ry := y % cs.sz
	rx := x % cs.sz
	if south > 0 {
		ny := ry - rx
		if ny%south == 0 {
			return color.Black
		}
	}
	if east > 0 {
		nx := rx - ry
		if nx%east == 0 {
			return color.Black
		}
	}
	return color.White
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (cs *ColourSource) Opaque() bool {
	return false
}
