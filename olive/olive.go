package olivego

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

// Fill fills the entire given canvas
func Fill(pixels [][]uint32, width, height int, color uint32) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y][x] = color
		}
	}
}

// FillRect renders a rectangle
func FillRect(pixels [][]uint32, pixelsWidth, pixelsHeight int, x0, y0, w, h int, color uint32) {
	for dy := 0; dy < h; dy++ {
		y := y0 + dy
		if y >= 0 && y < pixelsHeight {
			for dx := 0; dx < w; dx++ {
				x := x0 + dx
				if x >= 0 && x < pixelsWidth {
					pixels[y][x] = color
				}
			}
		}
	}
}

// SaveToPpm saves the pixels into the ppm format, generating a ppm file
func SaveToPpm(pixels [][]uint32, width, height int, filePath string) error {
	if filePath == "" {
		return errors.New("file path must be provided")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "P6\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]
			r := uint8(pixel >> 24) // red in RGBA
			g := uint8(pixel >> 16) // green in RGBA
			b := uint8(pixel >> 8)  // green in RGBA
			// ignore alpha channel
			_, err := file.Write([]byte{r, g, b})
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// SaveToPng saves the pixels into the png format, generating a png file
func SaveToPng(pixels [][]uint32, width, height int, filePath string) error {
	if filePath == "" {
		return errors.New("file path must be provided")
	}

	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]
			r := uint8(pixel >> 24)
			g := uint8(pixel >> 16)
			b := uint8(pixel >> 8)
			a := uint8(pixel)
			img.SetNRGBA(x, y, color.NRGBA{R: r, G: g, B: b, A: a})
		}
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = png.Encode(file, img)
	if err != nil {
		return err
	}

	return nil
}

// FillCircle renders a circle with given coordinates and radius
func FillCircle(pixels [][]uint32, pixelsWidth, pixelsHeight int, cx, cy, r int, color uint32) {
	x1 := cx - r
	x2 := cx + r
	y1 := cx - r
	y2 := cx + r

	for y := y1; y <= y2; y++ {
		if y >= 0 && pixelsHeight > y {
			for x := x1; x <= x2; x++ {
				if x >= 0 && pixelsWidth > x {
					dx := x - cx
					dy := y - cy
					if dx*dx+dy*dy <= r*r {
						pixels[y][x] = color
					}
				}
			}
		}
	}
}

// Line draws a straight line between 2 points: (x0, y0), (x1, y1)
func Line(pixels [][]uint32, pixelsWidth, pixelsHeight int, x0, y0, x1, y1 int, color uint32) {
	dx := x1 - x0
	dy := y1 - y0
	var m float64

	if dx != 0 {
		m = float64(dy) / float64(dx)
	}

	n := y0 - int(m*float64(x0))

	steep := abs(dx) < abs(dy)

	if steep {
		x0, y0 = y0, x0
		x1, y1 = y1, x1
	}

	if x0 > x1 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}

	if dx != 0 {
		m = float64(y1-y0) / float64(x1-x0)
	}
	n = y0 - int(m*float64(x0))

	for x := x0; x < x1; x++ {
		var y int
		if steep {
			y = int(m*float64(x) + float64(n))

			if y >= 0 && x < pixelsWidth && x >= 0 && y < pixelsHeight {
				pixels[y][x] = color
			}
		} else {
			y = int(m*float64(x) + float64(n))

			if x >= 0 && x < pixelsWidth && y >= 0 && y < pixelsHeight {
				pixels[y][x] = color
			}
		}
	}
}

// FillTriangle renders a triangle based on given 3 coordinates
func FillTriangle(pixels [][]uint32, pixelsWidth, pixelsHeight, x0, y0, x1, y1, x2, y2 int, color uint32) {
	if y1 < y0 {
		x0, x1 = x1, x0
		y0, y1 = y1, y0
	}
	if y2 < y0 {
		x0, x2 = x2, x0
		y0, y2 = y2, y0
	}
	if y2 < y1 {
		x1, x2 = x2, x1
		y1, y2 = y2, y1
	}

	interpolate := func(x0, y0, x1, y1, y int) int {
		if y0 == y1 {
			return x0
		}
		return x0 + (x1-x0)*(y-y0)/(y1-y0)
	}

	drawLine := func(y, x1, x2 int) {
		for x := x1; x <= x2; x++ {
			if x >= 0 && x < pixelsWidth && y >= 0 && y < pixelsHeight {
				pixels[y][x] = color
			}
		}
	}

	for y := y0; y <= y1; y++ {
		xa := interpolate(x0, y0, x1, y1, y)
		xb := interpolate(x0, y0, x2, y2, y)
		if xa > xb {
			xa, xb = xb, xa
		}
		drawLine(y, xa, xb)
	}

	for y := y1; y <= y2; y++ {
		xa := interpolate(x1, y1, x2, y2, y)
		xb := interpolate(x0, y0, x2, y2, y)
		if xa > xb {
			xa, xb = xb, xa
		}
		drawLine(y, xa, xb)
	}
}

// BlendColors blends two given colors
func BlendColors(bgColor, fgColor uint32) uint32 {
	bgR := uint8((bgColor >> 24) & 0xFF)
	bgG := uint8((bgColor >> 16) & 0xFF)
	bgB := uint8((bgColor >> 8) & 0xFF)

	fgR := uint8((fgColor >> 24) & 0xFF)
	fgG := uint8((fgColor >> 16) & 0xFF)
	fgB := uint8((fgColor >> 8) & 0xFF)
	fgA := uint8(fgColor & 0xFF)

	alpha := float64(fgA) / 255.0

	r := uint8(float64(bgR)*(1-alpha) + float64(fgR)*alpha)
	g := uint8(float64(bgG)*(1-alpha) + float64(fgG)*alpha)
	b := uint8(float64(bgB)*(1-alpha) + float64(fgB)*alpha)

	return (uint32(r) << 24) | (uint32(g) << 16) | (uint32(b) << 8) | 0xFF
}

// ApplyBlend applies the blend into the pixels
func ApplyBlend(pixels [][]uint32, x, y int, fgColor uint32) {
	if x < 0 || x >= len(pixels[0]) || y < 0 || y >= len(pixels) {
		return // out of bounds
	}

	bgColor := pixels[y][x]
	blendedColor := BlendColors(bgColor, fgColor)
	pixels[y][x] = blendedColor
}

// TODO: DrawCircle
// TODO: {Draw&&Fill}Elipse
