package main

import (
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/bmp"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"log"
	"os"
)

func main() {
	log.SetFlags(log.Flags() | log.Lshortfile)
	log.Printf("Setup")
	pal := []color.Color{
		color.White,
		color.Black,
	}
	totalSize := image.Rect(0,0,200,200)
	i := image.NewPaletted(totalSize, pal)
	log.Print("Adding header")
	fc, err := truetype.Parse(gomono.TTF)
	if err != nil {
		log.Panicf("Fond parse error: %#v", err)
	}
	face := truetype.NewFace(fc, &truetype.Options{
		Size: 16,
		DPI:  150,
	})
	headerEnd := face.Metrics().Ascent + face.Metrics().Descent
	d := &font.Drawer{
		Dst:  i,
		Src:  image.NewUniform(color.Black),
		Face: face,
		Dot:  fixed.Point26_6{
			X: fixed.I(0),
			Y: headerEnd,
		},
	}
	d.DrawString("Hi")
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
