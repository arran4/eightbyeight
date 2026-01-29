package eightbyeight

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"os"
	"testing"
)

// Custom BMP reader for the specific formats found in exampledata
func readBMP(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	// Read File Header (14 bytes) + Info Header (40 bytes) = 54 bytes
	header := make([]byte, 54)
	if _, err := io.ReadFull(f, header); err != nil {
		return nil, err
	}

	if header[0] != 'B' || header[1] != 'M' {
		return nil, fmt.Errorf("not a BMP file")
	}

	dataOffset := binary.LittleEndian.Uint32(header[10:14])
	width := int32(binary.LittleEndian.Uint32(header[18:22]))
	height := int32(binary.LittleEndian.Uint32(header[22:26]))
	bpp := binary.LittleEndian.Uint16(header[28:30])
	compression := binary.LittleEndian.Uint32(header[30:34])

	if compression != 0 {
		return nil, fmt.Errorf("unsupported compression: %d", compression)
	}

	var palette []color.Color
	if bpp <= 8 {
		paletteSize := int(dataOffset) - 54
		numColors := paletteSize / 4
		if numColors > 0 {
			pData := make([]byte, paletteSize)
			if _, err := io.ReadFull(f, pData); err != nil {
				return nil, err
			}
			for i := 0; i < numColors; i++ {
				// BMP palette is BGR(A)
				b := pData[i*4]
				g := pData[i*4+1]
				r := pData[i*4+2]
				palette = append(palette, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}

	// Seek to data offset
	if _, err := f.Seek(int64(dataOffset), 0); err != nil {
		return nil, err
	}

	img := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), palette)

	// Rows are stored bottom-up
	// Row size padded to 4 bytes
	var rowSize int
	if bpp == 1 {
		rowSize = (int(width) + 7) / 8
	} else if bpp == 4 {
		rowSize = (int(width) + 1) / 2
	} else {
		return nil, fmt.Errorf("unsupported bpp: %d", bpp)
	}
	padding := (4 - (rowSize % 4)) % 4
	stride := rowSize + padding

	rowData := make([]byte, stride)

	for y := int(height) - 1; y >= 0; y-- {
		if _, err := io.ReadFull(f, rowData); err != nil {
			return nil, err
		}
		for x := 0; x < int(width); x++ {
			var colorIdx uint8
			if bpp == 1 {
				// 1 bit per pixel. MSB first? Windows BMP usually MSB left-most pixel.
				byteIdx := x / 8
				bitIdx := 7 - (x % 8)
				colorIdx = (rowData[byteIdx] >> bitIdx) & 1
			} else if bpp == 4 {
				// 4 bits per pixel. High nibble is left pixel.
				byteIdx := x / 2
				if x%2 == 0 {
					colorIdx = (rowData[byteIdx] >> 4) & 0x0F
				} else {
					colorIdx = rowData[byteIdx] & 0x0F
				}
			}
			if int(colorIdx) < len(palette) {
				img.SetColorIndex(x, y, colorIdx)
			}
		}
	}

	return img, nil
}

func compareImages(img1, img2 image.Image) (int, float64, error) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	intersect := b1.Intersect(b2)
	if intersect.Empty() {
		return 0, 0, fmt.Errorf("no intersection between %v and %v", b1, b2)
	}

	diffPixels := 0
	totalPixels := intersect.Dx() * intersect.Dy()

	for y := intersect.Min.Y; y < intersect.Max.Y; y++ {
		for x := intersect.Min.X; x < intersect.Max.X; x++ {
			c1 := img1.At(x, y)
			c2 := img2.At(x, y)

			r1, g1, b1, _ := c1.RGBA()
			r2, g2, b2, _ := c2.RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 {
				diffPixels++
			}
		}
	}

	// Calculate mismatch area penalty?
	// For now, just return difference within intersection.
	return diffPixels, float64(diffPixels) / float64(totalPixels), nil
}

func TestReproduceExampleData(t *testing.T) {
	tests := []struct {
		filename            string
		expectedTitle       string
		expectedColors      []color.Color
		expectedRows        int
		expectedCols        int
		expectedFontSize    float64
		expectedDPI         float64
		expectedLabelSizing string
	}{
		{
			filename:      "exampledata/128BWGR.BMP",
			expectedTitle: "128 BWGR", // Guess
			// Palette: Black, White, Green, Red. But usually GridBuilder sets White then Black?
			// example 128BWGR uses 4 colors.
			expectedColors:      []color.Color{color.RGBA{0, 0, 0, 255}, color.RGBA{255, 255, 255, 255}, color.RGBA{0, 255, 0, 255}, color.RGBA{255, 0, 0, 255}},
			expectedRows:        5,
			expectedCols:        10,
			expectedFontSize:    8,
			expectedDPI:         150,
			expectedLabelSizing: "  255",
		},
		{
			filename:            "exampledata/COLRMODS.BMP",
			expectedTitle:       "Colour Modes",                          // Guess
			expectedColors:      []color.Color{color.White, color.Black}, // Placeholder
			expectedRows:        12,
			expectedCols:        16,
			expectedFontSize:    8,
			expectedDPI:         150,
			expectedLabelSizing: "  255",
		},
		{
			filename:            "exampledata/EARLYRED.BMP",
			expectedTitle:       "Early Red",                             // Guess
			expectedColors:      []color.Color{color.White, color.Black}, // Placeholder
			expectedRows:        12,
			expectedCols:        16,
			expectedFontSize:    8,
			expectedDPI:         150,
			expectedLabelSizing: "  255",
		},
	}

	for _, tc := range tests {
		t.Run(tc.filename, func(t *testing.T) {
			targetImg, err := readBMP(tc.filename)
			if err != nil {
				t.Fatalf("Failed to read BMP: %v", err)
			}

			b := NewGridBuilder().
				WithTitle(tc.expectedTitle).
				WithDimensions(tc.expectedRows, tc.expectedCols).
				WithColors(tc.expectedColors).
				WithFont(tc.expectedFontSize, tc.expectedDPI).
				WithLabelSizing(tc.expectedLabelSizing)

			genImg := b.Generate()

			diff, diffPct, err := compareImages(genImg, targetImg)
			if err != nil {
				t.Errorf("Comparison failed: %v", err)
			} else {
				if genImg.Bounds() != targetImg.Bounds() {
					t.Logf("Warning: bounds mismatch. Generated %v, Target %v. Comparing intersection.", genImg.Bounds(), targetImg.Bounds())
				}
				t.Logf("Difference in intersection: %d pixels (%.2f%%)", diff, diffPct*100)
			}
		})
	}
}
