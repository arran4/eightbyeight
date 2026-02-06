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
	"sort"
	"strings"
	"testing"
)

type Config struct {
	Title       string
	Colors      []color.Color
	LabelSizing string
	IDSequence  []int
}

// readBMP handles 1-bit and 4-bit BMPs
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

func isWhite(c color.Color) bool {
	r, g, b, _ := c.RGBA()
	return r > 0xF000 && g > 0xF000 && b > 0xF000
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

	configs := map[string]Config{
		"128BWGR.BMP": {
			Title:       "128 BWGR",
			LabelSizing: "  255",
			Colors: []color.Color{
				color.RGBA{0, 0, 0, 255},
				color.RGBA{255, 255, 255, 255},
			},
			// Reading order sequence
			IDSequence: []int{
				114, 110, 106, 102, 98, 14, 4,
				94, 90, 86, 82, 78, 10, 3,
				74, 70, 66, 62, 58, 9, 2,
				54, 50, 46, 42, 38, 8, 1,
				34, 30, 26, 22, 18, 7, 0,
				6, 5,
			},
		},
		"COLRMODS.BMP": {
			Title:       "Colour Modes",
			LabelSizing: "  255",
			Colors:      cgaPalette,
			// Assuming row-major 0..N
			IDSequence: makeRange(0, 255),
		},
		"EARLYRED.BMP": {
			Title:       "Early Red",
			LabelSizing: "  255",
			Colors:      cgaPalette,
			IDSequence:  makeRange(0, 255),
		},
	}

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

			// Segmentation
			blobs := findPatternBlobs(fullImg)
			if len(blobs) == 0 {
				t.Errorf("No patterns found in %s", filename)
				return
			}

			// Sort Blobs
			sortBlobsSpatial(blobs)

			// Process
			for i, r := range blobs {
				if i >= len(cfg.IDSequence) {
					break // Found more blobs than expected IDs
				}
				mode := cfg.IDSequence[i]

				// Extract
				subImg := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
				draw.Draw(subImg, subImg.Bounds(), fullImg, r.Min, draw.Src)

				// Save
				outFile := filepath.Join(outDir, fmt.Sprintf("%d.png", mode))
				f, err := os.Create(outFile)
				if err != nil {
					t.Errorf("Failed to create %s: %v", outFile, err)
					continue
				}
				png.Encode(f, subImg)
				f.Close()

				// Verify
				t.Run(fmt.Sprintf("Mode_%d", mode), func(t *testing.T) {
					src := NewColourSource(mode, cfg.Colors...)
					genImg := image.NewRGBA(image.Rect(0, 0, r.Dx(), r.Dy()))
					draw.Draw(genImg, genImg.Bounds(), src, image.Point{}, draw.Src)

					diff, diffPct, err := compareImages(genImg, subImg)
					if err != nil {
						t.Errorf("Comparison failed: %v", err)
					} else if diff > 0 {
						t.Errorf("Difference: %d pixels (%.2f%%)", diff, diffPct*100)
					}
				})
			}
		})
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func findPatternBlobs(img image.Image) []image.Rectangle {
	bounds := img.Bounds()
	visited := make([]bool, bounds.Dx()*bounds.Dy())
	var blobs []image.Rectangle

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			idx := (y-bounds.Min.Y)*bounds.Dx() + (x - bounds.Min.X)
			if visited[idx] {
				continue
			}

			if !isWhite(img.At(x, y)) {
				// BFS
				minX, maxX := x, x
				minY, maxY := y, y
				q := []image.Point{{x, y}}
				visited[idx] = true

				for len(q) > 0 {
					p := q[0]
					q = q[1:]

					if p.X < minX {
						minX = p.X
					}
					if p.X > maxX {
						maxX = p.X
					}
					if p.Y < minY {
						minY = p.Y
					}
					if p.Y > maxY {
						maxY = p.Y
					}

					dirs := []image.Point{{0, 1}, {0, -1}, {1, 0}, {-1, 0}}
					for _, d := range dirs {
						nx, ny := p.X+d.X, p.Y+d.Y
						if nx >= bounds.Min.X && nx < bounds.Max.X && ny >= bounds.Min.Y && ny < bounds.Max.Y {
							nIdx := (ny-bounds.Min.Y)*bounds.Dx() + (nx - bounds.Min.X)
							if !visited[nIdx] {
								if !isWhite(img.At(nx, ny)) {
									visited[nIdx] = true
									q = append(q, image.Point{nx, ny})
								}
							}
						}
					}
				}

				rect := image.Rect(minX, minY, maxX+1, maxY+1)
				// Filter for "Pattern" size (heuristic: patterns are larger than labels)
				// 110.png is 67x39. Labels are small text.
				// Let's filter > 20x20
				if rect.Dx() > 20 && rect.Dy() > 20 {
					blobs = append(blobs, rect)
				}
			}
		}
	}
	return blobs
}

func sortBlobsSpatial(blobs []image.Rectangle) {
	sort.Slice(blobs, func(i, j int) bool {
		// Row binning logic
		// Assume rows are separated by at least 20px
		if blobs[i].Min.Y/40 != blobs[j].Min.Y/40 {
			return blobs[i].Min.Y < blobs[j].Min.Y
		}
		return blobs[i].Min.X < blobs[j].Min.X
	})
}
