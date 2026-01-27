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
	"math"
	"os"
)

type GridBuilder struct {
	Title    string
	Rows     int
	Columns  int
	CellSize int
	Palette  []color.Color
}

func NewGridBuilder() *GridBuilder {
	return &GridBuilder{
		Title:    "Grid Draw",
		Rows:     10,
		Columns:  4,
		CellSize: 64,
		Palette:  []color.Color{color.White, color.Black},
	}
}

func (b *GridBuilder) WithTitle(title string) *GridBuilder {
	b.Title = title
	return b
}

func (b *GridBuilder) WithDimensions(rows, columns int) *GridBuilder {
	b.Rows = rows
	b.Columns = columns
	return b
}

func (b *GridBuilder) WithColors(palette []color.Color) *GridBuilder {
	b.Palette = palette
	return b
}

func (b *GridBuilder) Generate() image.Image {
	log.Printf("Setup")
	fc, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Panicf("Font parse error: %#v", err)
	}

	lines := b.Rows
	lineLength := b.Columns

	fontFace := truetype.NewFace(fc, &truetype.Options{
		Size: 16,
		DPI:  150,
	})
	fontHeight := fontFace.Metrics().Ascent
	lineHeight := fontFace.Metrics().Height + fontFace.Metrics().Descent

	labelBounds, _ := font.BoundString(fontFace, fmt.Sprintf("__%d__", 255))
	titleBounds, _ := font.BoundString(fontFace, b.Title)

	totalSize := image.Rect(0, 0, IntMax(titleBounds.Max.X.Ceil(), IntMax(labelBounds.Max.X.Ceil(), b.CellSize)*lineLength), (lineHeight.Ceil()+b.CellSize)*(lines)+lineHeight.Ceil())

	// Use the first color in palette as background if available, otherwise White
	bg := color.Color(color.White)
	if len(b.Palette) > 0 {
		bg = b.Palette[0]
	}

	i := image.NewPaletted(totalSize, b.Palette)
	// Fill with background
	draw.Draw(i, i.Bounds(), image.NewUniform(bg), image.Point{}, draw.Src)

	log.Print("Adding header")
	// Text color: use Black if palette has it, or just Black.
	// Actually, for paletted image, we should use a color from palette.
	// The original code used a palette of White, Black.
	// And drew text with Src = image.NewUniform(color.Black).

	textColor := color.Color(color.Black)
	// If palette has > 1 color, maybe use the second one as text color?
	if len(b.Palette) > 1 {
		textColor = b.Palette[1]
	}

	d := &font.Drawer{
		Dst:  i,
		Src:  image.NewUniform(textColor),
		Face: fontFace,
		Dot: fixed.Point26_6{
			X: fixed.I(0),
			Y: fontHeight,
		},
	}
	d.DrawString(b.Title)
	log.Printf("Drawing grid with labels")
	for y := 0; y < lines; y++ {
		yTop := lineHeight.Ceil() + (lineHeight.Ceil()+b.CellSize)*(y)
		for x := 0; x < lineLength; x++ {
			xLeft := IntMax(labelBounds.Max.X.Ceil(), b.CellSize) * x
			r := image.Rect(
				xLeft,
				yTop,
				IntMax(labelBounds.Max.X.Ceil(), b.CellSize)*(x+1)-1,
				yTop+b.CellSize-1,
			)
			d.Dot.X = fixed.I(IntMax(labelBounds.Max.X.Ceil(), b.CellSize) * (x))
			d.Dot.Y = fixed.I(yTop + b.CellSize - 1 + fontHeight.Ceil())
			dy, dx := r.Dy(), r.Dx()
			if dy > b.CellSize {
				r.Min.Y += (dy - b.CellSize) / 2
				r.Max.Y -= (dy - b.CellSize) / 2
			}
			if dx > b.CellSize {
				r.Min.X += (dx - b.CellSize) / 2
				r.Max.X -= (dx - b.CellSize) / 2
			}
			mode := x + y*lineLength
			d.DrawString(fmt.Sprintf("  %d", mode))
			// Pass the palette to NewColourSource
			draw.Draw(i, r, NewColourSource(mode, b.Palette...), image.Point{}, draw.Src)
		}
	}
	return i
}

func (b *GridBuilder) Save(filename string) error {
	img := b.Generate()
	log.Print("Writing file")
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer f.Close()
	if err := bmp.Encode(f, img); err != nil {
		return err
	}
	log.Printf("Done")
	return nil
}

func IntMax(i1, i2 int) int {
	if i1 > i2 {
		return i1
	}
	return i2
}

func NewColourSource(mode int, colors ...color.Color) image.Image {
	sz := 8
	south := [4]int{
		1,
		0,
		0,
		0,
	}
	n := mode
	for i := 0; i < 4 && n > 0; i++ {
		south[i] = n % 4
		n /= 4
	}
	return &ColourSource{
		mode:   mode,
		colors: colors,
		sz:     sz,
		south:  south,
	}
}

type ColourSource struct {
	colors []color.Color
	mode   int
	sz     int
	south  [4]int
}

func (cs *ColourSource) ColorModel() color.Model {
	return color.ModelFunc(func(c color.Color) color.Color { return c })
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
	xp := x % cs.sz
	dp := (y + cs.sz - xp) % cs.sz

	if cs.mode >= xp {
		xp = 3 - (3 - xp)
		if xp < 0 {
			xp = -xp
		}
		sv := cs.south[xp%len(cs.south)]
		sv = int(math.Pow(float64(2), float64(4-sv)))
		if sv > 0 && dp%sv == 0 {
			// Use the second color (index 1) for "foreground"/black if available
			if len(cs.colors) > 1 {
				return cs.colors[1]
			}
			return color.Black
		}
	}

	// Use the first color (index 0) for "background"/white if available
	if len(cs.colors) > 0 {
		return cs.colors[0]
	}
	return color.White
}

// Opaque scans the entire image and reports whether it is fully opaque.
func (cs *ColourSource) Opaque() bool {
	return false
}
