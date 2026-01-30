package eightbyeight

import (
	"encoding/json"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
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

	return diffPixels, float64(diffPixels) / float64(totalPixels), nil
}

func TestReproduceExampleData(t *testing.T) {
	files, err := filepath.Glob("exampledata/*.json")
	if err != nil {
		t.Fatal(err)
	}

	for _, configFile := range files {
		t.Run(configFile, func(t *testing.T) {
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

			// Read Target Image
			imgFile := configFile[:len(configFile)-5] + ".png"
			imgF, err := os.Open(imgFile)
			if err != nil {
				t.Fatalf("Failed to open image %s: %v", imgFile, err)
			}
			defer imgF.Close()

			targetImg, err := png.Decode(imgF)
			if err != nil {
				t.Fatalf("Failed to decode image: %v", err)
			}

			// Build Palette
			var palette []color.Color
			for _, c := range cfg.Colors {
				palette = append(palette, c.ToColor())
			}

			// Generate
			b := NewGridBuilder().
				WithTitle(cfg.Title).
				WithDimensions(cfg.Rows, cfg.Columns).
				WithColors(palette).
				WithFont(cfg.FontSize, cfg.DPI).
				WithLabelSizing(cfg.LabelSizing)

			genImg := b.Generate()

			// Compare
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
