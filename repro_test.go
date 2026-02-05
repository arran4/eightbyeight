package eightbyeight

import (
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
	"golang.org/x/image/font/gofont/gomono"
)

type Config struct {
	Title       string
	Rows        int
	Columns     int
	Colors      []color.Color
	FontSize    float64
	DPI         float64
	LabelSizing string
	// CustomGrid allows defining explicit cells for irregular layouts
	CustomGrid []CustomCell
}

type CustomCell struct {
	Mode int
	Rect image.Rectangle
}

// readBMP (kept from previous steps)
func readBMP(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

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
				b := pData[i*4]
				g := pData[i*4+1]
				r := pData[i*4+2]
				palette = append(palette, color.RGBA{R: r, G: g, B: b, A: 255})
			}
		}
	}

	if _, err := f.Seek(int64(dataOffset), 0); err != nil {
		return nil, err
	}

	img := image.NewPaletted(image.Rect(0, 0, int(width), int(height)), palette)

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
				byteIdx := x / 8
				bitIdx := 7 - (x % 8)
				colorIdx = (rowData[byteIdx] >> bitIdx) & 1
			} else if bpp == 4 {
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

func isUniform(img image.Image) bool {
	b := img.Bounds()
	if b.Dx() == 0 || b.Dy() == 0 {
		return true
	}
	first := img.At(b.Min.X, b.Min.Y)
	r1, g1, b1, a1 := first.RGBA()

	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := img.At(x, y)
			r, g, b, a := c.RGBA()
			if r != r1 || g != g1 || b != b1 || a != a1 {
				return false
			}
		}
	}
	return true
}

func compareImages(img1, img2 image.Image) (int, float64, error) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	// Strict bounds check
	if b1.Dx() != b2.Dx() || b1.Dy() != b2.Dy() {
		return 0, 0, fmt.Errorf("bounds mismatch: %v != %v", b1, b2)
	}

	diffPixels := 0
	totalPixels := b1.Dx() * b1.Dy()

	for y := 0; y < b1.Dy(); y++ {
		for x := 0; x < b1.Dx(); x++ {
			c1 := img1.At(b1.Min.X+x, b1.Min.Y+y)
			c2 := img2.At(b2.Min.X+x, b2.Min.Y+y)

			r1, g1, b1, _ := c1.RGBA()
			r2, g2, b2, _ := c2.RGBA()

			if r1 != r2 || g1 != g2 || b1 != b2 {
				diffPixels++
			}
		}
	}

	return diffPixels, float64(diffPixels) / float64(totalPixels), nil
}

func TestReproducePatterns(t *testing.T) {
	// 1. Define Configs
	cgaPalette := []color.Color{
		color.RGBA{0, 0, 0, 255},       // 0: Black
		color.RGBA{128, 0, 0, 255},     // 1: Maroon
		color.RGBA{0, 128, 0, 255},     // 2: Green
		color.RGBA{128, 128, 0, 255},   // 3: Olive
		color.RGBA{0, 0, 128, 255},     // 4: Navy
		color.RGBA{128, 0, 128, 255},   // 5: Purple
		color.RGBA{0, 128, 128, 255},   // 6: Teal
		color.RGBA{128, 128, 128, 255}, // 7: Silver
		color.RGBA{192, 192, 192, 255}, // 8: Gray
		color.RGBA{255, 0, 0, 255},     // 9: Red
		color.RGBA{0, 255, 0, 255},     // 10: Lime
		color.RGBA{255, 255, 0, 255},   // 11: Yellow
		color.RGBA{0, 0, 255, 255},     // 12: Blue
		color.RGBA{255, 0, 255, 255},   // 13: Fuchsia
		color.RGBA{0, 255, 255, 255},   // 14: Aqua
		color.RGBA{255, 255, 255, 255}, // 15: White
	}

	// Generate custom grid for 128BWGR based on analysis
	// 110 is at (116, 112). Col stride 96. Row stride 80 (approx). Size 67x39.
	// Labels logic: Col 1 starts at 114, dec 20. Col 2 starts at 110, dec 20.
	customGrid128 := []CustomCell{}

	// Column X starts
	colXs := []int{20, 116, 212, 308, 404}
	// Row Y starts
	rowYs := []int{112, 192, 272, 352, 432}
	// Start Labels per column
	startLabels := []int{114, 110, 106, 102, 98}

	for cIdx, startX := range colXs {
		currentLabel := startLabels[cIdx]
		for _, startY := range rowYs {
			customGrid128 = append(customGrid128, CustomCell{
				Mode: currentLabel,
				Rect: image.Rect(startX, startY, startX+67, startY+39),
			})
			currentLabel -= 20
		}
	}

	configs := map[string]Config{
		"128BWGR.BMP": {
			Title: "128 BWGR",
			// Rows/Cols unused if CustomGrid present
			CustomGrid: customGrid128,
			Colors: []color.Color{
				color.RGBA{0, 0, 0, 255},
				color.RGBA{255, 255, 255, 255},
			},
		},
		"COLRMODS.BMP": {
			Title:       "Colour Modes",
			Rows:        16,
			Columns:     16,
			FontSize:    8,
			DPI:         150,
			LabelSizing: "  255",
			Colors:      cgaPalette,
		},
		"EARLYRED.BMP": {
			Title:       "Early Red",
			Rows:        16,
			Columns:     16,
			FontSize:    8,
			DPI:         150,
			LabelSizing: "  255",
			Colors:      cgaPalette,
		},
	}

	// 2. Setup Font
	fc, err := truetype.Parse(gomono.TTF)
	if err != nil {
		t.Fatalf("Font parse error: %#v", err)
	}

	// 3. Iterate
	for filename, cfg := range configs {
		t.Run(filename, func(t *testing.T) {
			filePath := filepath.Join("exampledata", filename)
			if _, err := os.Stat(filePath); os.IsNotExist(err) {
				t.Fatalf("Reference BMP missing: %s", filePath)
			}

			fullImg, err := readBMP(filePath)
			if err != nil {
				t.Fatalf("Failed to read BMP %s: %v", filePath, err)
			}

			datasetName := strings.TrimSuffix(filename, ".BMP")
			outDir := filepath.Join("exampledata", datasetName)
			if err := os.MkdirAll(outDir, 0755); err != nil {
				t.Fatalf("Failed to create dir %s: %v", outDir, err)
			}

			// Extraction loop
			if len(cfg.CustomGrid) > 0 {
				for _, cell := range cfg.CustomGrid {
					processCell(t, fullImg, cell.Rect, cell.Mode, cfg.Colors, outDir)
				}
			} else {
				// Standard Grid
				cellSize := 64 // default
				fontFace := truetype.NewFace(fc, &truetype.Options{
					Size: cfg.FontSize,
					DPI:  cfg.DPI,
				})
				lineHeight := fontFace.Metrics().Height + fontFace.Metrics().Descent
				labelBounds, _ := font.BoundString(fontFace, cfg.LabelSizing)

				for y := 0; y < cfg.Rows; y++ {
					yTop := lineHeight.Ceil() + (lineHeight.Ceil()+cellSize)*(y)
					for x := 0; x < cfg.Columns; x++ {
						xLeft := IntMax(labelBounds.Max.X.Ceil(), cellSize) * x
						width := IntMax(labelBounds.Max.X.Ceil(), cellSize)
						r := image.Rect(
							xLeft,
							yTop,
							xLeft+width-1,
							yTop+cellSize-1,
						)

						// Centering logic (from Builder)
						dy, dx := r.Dy(), r.Dx()
						if dy > cellSize {
							r.Min.Y += (dy - cellSize) / 2
							r.Max.Y -= (dy - cellSize) / 2
						}
						if dx > cellSize {
							r.Min.X += (dx - cellSize) / 2
							r.Max.X -= (dx - cellSize) / 2
						}

						mode := x + y*cfg.Columns
						processCell(t, fullImg, r, mode, cfg.Colors, outDir)
					}
				}
			}
		})
	}
}

func processCell(t *testing.T, fullImg image.Image, r image.Rectangle, mode int, palette []color.Color, outDir string) {
	intersect := r.Intersect(fullImg.Bounds())
	if intersect.Empty() {
		return
	}

	// Extract full cell/pattern rect
	subImg := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
	draw.Draw(subImg, subImg.Bounds(), fullImg, r.Min, draw.Src)

	if isUniform(subImg) {
		return
	}

	// Save reference PNG (Decimal Name)
	outFile := filepath.Join(outDir, fmt.Sprintf("%d.png", mode))
	f, err := os.Create(outFile)
	if err != nil {
		t.Errorf("Failed to create %s: %v", outFile, err)
		return
	}
	png.Encode(f, subImg)
	f.Close()

	// Verify
	t.Run(fmt.Sprintf("Mode_%d", mode), func(t *testing.T) {
		// Generate
		src := NewColourSource(mode, palette...)
		genImg := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
		draw.Draw(genImg, genImg.Bounds(), src, image.Point{}, draw.Src)

		// Compare
		diff, diffPct, err := compareImages(genImg, subImg)
		if err != nil {
			t.Errorf("Comparison failed: %v", err)
		} else if diff > 0 {
			t.Errorf("Difference: %d pixels (%.2f%%)", diff, diffPct*100)
		}
	})
}
