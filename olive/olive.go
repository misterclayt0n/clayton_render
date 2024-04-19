package olivego

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

func BuildPixel(height, width int) [][]uint32 {
	pixels := make([][]uint32, height)
	for i := range pixels {
		pixels[i] = make([]uint32, width)
	}

	return pixels
}

// Fill fills the entire given canvas
func Fill(pixels [][]uint32, width, height int, color uint32) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			blendPixel(pixels, x, y, color)
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
					blendPixel(pixels, x, y, color)
				}
			}
		}
	}
}

// DrawRect draws the outline of a rectangle.
func DrawRect(pixels [][]uint32, width, height, x0, y0, w, h int, color uint32) {
	// Top and bottom
	for x := x0; x < x0+w; x++ {
		blendPixel(pixels, x, y0, color)
		blendPixel(pixels, x, y0+h-1, color)
	}
	// Left and right
	for y := y0; y < y0+h; y++ {
		blendPixel(pixels, x0, y, color)
		blendPixel(pixels, x0+w-1, y, color)
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
	y1 := cy - r
	y2 := cy + r

	for y := y1; y <= y2; y++ {
		if y >= 0 && pixelsHeight > y {
			for x := x1; x <= x2; x++ {
				if x >= 0 && pixelsWidth > x {
					dx := x - cx
					dy := y - cy
					if dx*dx+dy*dy <= r*r {
						blendPixel(pixels, x, y, color)
					}
				}
			}
		}
	}
}

// DrawCircle draws the outline of a circle.
func DrawCircle(pixels [][]uint32, width, height, cx, cy, r int, color uint32) {
	x, y := r, 0
	p := 1 - r

	setCirclePixels(pixels, cx, cy, x, y, color)

	for x > y {
		y++
		if p <= 0 {
			p += 2*y + 1
		} else {
			x--
			p += 2*(y-x) + 1
		}
		setCirclePixels(pixels, cx, cy, x, y, color)
	}
}

// setCirclePixels sets the pixels for the 8 octants of the circle.
func setCirclePixels(pixels [][]uint32, cx, cy, x, y int, color uint32) {
	blendPixel(pixels, cx+x, cy+y, color)
	blendPixel(pixels, cx-x, cy+y, color)
	blendPixel(pixels, cx+x, cy-y, color)
	blendPixel(pixels, cx-x, cy-y, color)
	blendPixel(pixels, cx+y, cy+x, color)
	blendPixel(pixels, cx-y, cy+x, color)
	blendPixel(pixels, cx+y, cy-x, color)
	blendPixel(pixels, cx-y, cy-x, color)
}

// Line draws a straight line between 2 points: (x0, y0), (x1, y1)
func Line(pixels [][]uint32, pixelsWidth, pixelsHeight int, x0, y0, x1, y1 int, color uint32) {
	dx := abs(x1 - x0)
	dy := -abs(y1 - y0)
	sx := -1
	sy := -1
	if x0 < x1 {
		sx = 1
	}
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		blendPixel(pixels, x0, y0, color)

		if x0 == x1 && y0 == y1 {
			break
		}

		e2 := 2 * err
		if e2 >= dy {
			if x0 == x1 {
				break
			}
			err += dy
			x0 += sx
		}
		if e2 <= dx {
			if y0 == y1 {
				break
			}
			err += dx
			y0 += sy
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
				blendPixel(pixels, x, y, color)
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

// DrawTriangle draws the outline of a triangle.
func DrawTriangle(pixels [][]uint32, width, height, x0, y0, x1, y1, x2, y2 int, color uint32) {
	Line(pixels, width, height, x0, y0, x1, y1, color)
	Line(pixels, width, height, x1, y1, x2, y2, color)
	Line(pixels, width, height, x2, y2, x0, y0, color)
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

// blendPixel blends the foreground color with the background color directly on the canvas.
func blendPixel(pixels [][]uint32, x, y int, fgColor uint32) {
	if x < 0 || x >= len(pixels[0]) || y < 0 || y >= len(pixels) {
		return // out of bounds
	}
	bgColor := pixels[y][x]
	blendedColor := BlendColors(bgColor, fgColor)
	pixels[y][x] = blendedColor
}

// PixelsToBytes converts a bidimensional uint32 slice to a slice of bytes.
func PixelsToBytes(pixels [][]uint32) []byte {
	height := len(pixels)
	if height == 0 {
		return nil
	}
	width := len(pixels[0])
	bytes := make([]byte, width*height*4)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixel := pixels[y][x]
			i := (y*width + x) * 4
			bytes[i] = byte(pixel >> 24)
			bytes[i+1] = byte(pixel>>16) & 0xFF
			bytes[i+2] = byte(pixel>>8) & 0xFF
			bytes[i+3] = byte(pixel) & 0xFF
		}
	}

	return bytes
}

func RotatePoint(cx, cy, x, y, angle float64) (float64, float64) {
	rad := degreesToRadians(angle)
	cos := math.Cos(rad)
	sin := math.Sin(rad)

	x -= cx
	y -= cy

	newX := x*cos - y*sin
	newY := x*sin + y*cos

	newX += cx
	newY += cy
	return newX, newY
}

// TODO: DrawCircle
// TODO: {Draw&&Fill}Elipse
