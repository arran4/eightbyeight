package eightbyeight

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
)

type Config struct {
	Title       string        `json:"title"`
	Rows        int           `json:"rows"`
	Columns     int           `json:"columns"`
	Colors      []ColorConfig `json:"colors"`
	FontSize    float64       `json:"fontSize"`
	DPI         float64       `json:"dpi"`
	LabelSizing string        `json:"labelSizing"`
}

type ColorConfig struct {
	R, G, B, A uint8
}

func (c ColorConfig) ToColor() color.Color {
	return color.RGBA{c.R, c.G, c.B, c.A}
}

func compareImages(img1, img2 image.Image) (int, float64, error) {
	b1 := img1.Bounds()
	b2 := img2.Bounds()

	// Strict bounds check for 8x8 cells
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
	configFiles, err := filepath.Glob("exampledata/*.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, configFile := range configFiles {
		t.Run(configFile, func(t *testing.T) {
			datasetName := strings.TrimSuffix(filepath.Base(configFile), ".json")

			// Read Config
			f, err := os.Open(configFile)
			if err != nil {
				t.Fatalf("Failed to open config: %v", err)
			}
			defer f.Close()

			var cfg Config
			if err := json.NewDecoder(f).Decode(&cfg); err != nil {
				t.Fatalf("Failed to decode config: %v", err)
			}

			// Build Palette
			var palette []color.Color
			for _, c := range cfg.Colors {
				palette = append(palette, c.ToColor())
			}

			// Iterate through expected sub-images
			cellsDir := filepath.Join("exampledata", datasetName)
			cellFiles, err := filepath.Glob(filepath.Join(cellsDir, "*.png"))
			if err != nil {
				t.Fatalf("Failed to glob cells: %v", err)
			}
			if len(cellFiles) == 0 {
				t.Fatalf("No cell images found in %s", cellsDir)
			}

			for _, cellFile := range cellFiles {
				// Filename is mode.png
				base := filepath.Base(cellFile)
				modeStr := strings.TrimSuffix(base, ".png")
				mode, err := strconv.Atoi(modeStr)
				if err != nil {
					t.Logf("Skipping non-numeric file: %s", cellFile)
					continue
				}

				t.Run(fmt.Sprintf("Mode_%d", mode), func(t *testing.T) {
					// Read Target Image
					imgF, err := os.Open(cellFile)
					if err != nil {
						t.Fatalf("Failed to open image %s: %v", cellFile, err)
					}
					defer imgF.Close()

					targetImg, err := png.Decode(imgF)
					if err != nil {
						t.Fatalf("Failed to decode image: %v", err)
					}

					// Generate 8x8 pattern
					src := NewColourSource(mode, palette...)
					// ColourSource is infinite, need to draw it into an 8x8 image
					genImg := image.NewRGBA(image.Rect(0, 0, 8, 8))
					draw.Draw(genImg, genImg.Bounds(), src, image.Point{}, draw.Src)

					// Compare
					diff, diffPct, err := compareImages(genImg, targetImg)
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
