# EightByEight Pattern Generator

This project generates PNG images containing grids of 8x8 pixel patterns. It explores different permutations of patterns based on a modular arithmetic logic for colour dithering.

## Description

The tool creates a grid where each cell contains a unique 8x8 pattern. It renders a title and labels for the grid. The patterns are generated using a "ColourSource" algorithm that determines the color of each pixel based on its coordinates and a mode index.

## Usage

To run the project, you need Go installed.

```bash
go run cmd/eightbyeight/main.go
```

This will generate four files in the current directory:
- `out_bw.png`: Classic Black on White
- `out_terminal.png`: Green on Black (Terminal style)
- `out_solarized.png`: Solarized Light color scheme
- `out_mixing.png`: CGA Color Mixing

## Output

The output are PNG images with a title, a grid of patterns, and labels.

## Examples

Examples of generated patterns can be found in the `exampledata/` directory.

Here are the generated outputs:

### Classic - Black on White
![Classic Output](out_bw.png)

### Terminal - Green on Black
![Terminal Output](out_terminal.png)

### Solarized Light
![Solarized Output](out_solarized.png)

### CGA Color Mixing
![CGA Mixing Output](out_mixing.png)

## Builder

The project includes a `GridBuilder` to programmatically configure and generate these pattern grids.

```go
import "github.com/arran4/eightbyeight"

// ...

builder := eightbyeight.NewGridBuilder().
    WithTitle("My Custom Grid").
    WithDimensions(10, 5). // 10 rows, 5 columns
    WithColors([]color.Color{color.White, color.RGBA{255, 0, 0, 255}})

img := builder.Generate()
builder.Save("my_grid.bmp")
```
