package chip

import "image/color"

const PIXEL_SIZE = float32(10)

type Screen struct {
	Width  int
	Height int
	Pixels [][]color.Color
}

func NewScreen(width int, height int) *Screen {
	pixels := [][]color.Color{}
	for r := 0; r < height; r++ {
		row := []color.Color{}
		for c := 0; c < width; c++ {
			row = append(row, color.Black)
		}
		pixels = append(pixels, row)
	}
	return &Screen{
		Width:  width,
		Height: height,
		Pixels: pixels,
	}
}

func (screen *Screen) IsPixelOn(row int, col int) bool {
	return screen.Pixels[row][col] == color.White
}

func (screen *Screen) SetPixel(row int, col int, on bool) {
	c := color.Black
	if on {
		c = color.White
	}
	screen.Pixels[row][col] = c
}

func (screen *Screen) clear() {
	for _, row := range screen.Pixels {
		for i := range row {
			row[i] = color.Black
		}
	}
}

func (screen *Screen) UpdateScreen() {
	// for _, row := range screen.Pixels {
	// 	for i, pix := range row {
	// 		r, g, b, a := pix.RGBA()
	// 		r = r / 257
	// 		g = g / 257
	// 		b = b / 257
	// 		a = a / 257
	// 		row[i] = color.RGBA{
	// 			R: uint8(float32(r) * .999),
	// 			G: uint8(float32(g) / 2),
	// 			B: uint8(float32(b) / 2),
	// 			A: uint8(a),
	// 		}
	// 	}
	// }
}
