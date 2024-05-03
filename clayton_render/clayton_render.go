package clayton_render

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"math"
	"os"
)

// Canvas is the primary structure of clayton_render. Use it to begin the rendering.
type Canvas struct {
	Pixels []uint32
	Height int
	Width  int
	Stride int
}

// NormalizedRect contains the normalized values of a rectangle
type normalizedRect struct {
	x1, y1, x2, y2 int
}

func (c *Canvas) NormalizeRect(x, y, w, h int) (normalizedRect, bool) {
	// check if x is out of bounds
	if x < 0 {
		w += x
		x = 0
	}

	if x+w > c.Width {
		w = c.Width - x
	}

	// check if y is out of bounds
	if y < 0 {
		h += y
		y = 0
	}

	if y+h > c.Height {
		h = c.Height - y
	}

	if w <= 0 || h <= 0 {
		return normalizedRect{}, false
	}

	return normalizedRect{
		x1: x,
		y1: y,
		x2: x + w,
		y2: y + h,
	}, true
}

// NewCanvas builds the canvas with given width and height.
func NewCanvas(width, height int) *Canvas {
	pixels := make([]uint32, width*height)
	return &Canvas{
		Pixels: pixels,
		Width:  width,
		Height: height,
		Stride: width,
	}
}

// SetPixel set a given pixel to some color.
func (c *Canvas) SetPixel(x, y int, color uint32) {
	if x < 0 || x >= c.Width || y < 0 || y >= c.Height {
		return // out of bounds
	}

	index := y*c.Stride + x
	bgColor := c.Pixels[index]
	blendedColor := BlendColors(bgColor, color)
	c.Pixels[index] = blendedColor
}

// Fill fills the entire canvas with a desired color
func (c *Canvas) Fill(color uint32) {
	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			c.SetPixel(x, y, color)
		}
	}
}

// FillRect renders a rectangle
func (c *Canvas) FillRect(x0, y0, w, h int, color uint32) {
	nr, visible := c.NormalizeRect(x0, y0, w, h)

	if !visible {
		return // out of bounds
	}

	for dy := nr.y1; dy <= nr.y2; dy++ {
		for dx := nr.x1; dx <= nr.x2; dx++ {
			c.SetPixel(dx, dy, color)
		}
	}
}

// DrawRect draws the outline of a rectangle.
func (c *Canvas) DrawRect(x0, y0, w, h int, color uint32) {
	// top and bottom
	for x := x0; x < x0+w; x++ {
		c.SetPixel(x, y0, color)
		c.SetPixel(x, y0+h-1, color)
	}

	for y := y0; y < y0+h; y++ {
		c.SetPixel(x0, y, color)
		c.SetPixel(x0+w-1, y, color)
	}
}

// FillCircle renders a circle with the given coordinates and radius.
func (c *Canvas) FillCircle(cx, cy, r int, color uint32) {
	for y := cy - r - 1; y <= cy+r+1; y++ {
		for x := cx - r - 1; x <= cx+r+1; x++ {
			dx := x - cx
			dy := y - cy
			distanceSquared := dx*dx + dy*dy

			distanceToEdge := math.Sqrt(float64(distanceSquared)) - float64(r)

			if distanceToEdge < 0 {
				c.SetPixel(x, y, color)
			} else if distanceToEdge < 1 {
				opacity := 1.0 - distanceToEdge

				newAlpha := uint8(float64(color&0xFF) * opacity)
				newColor := (color & 0xFFFFFF00) | uint32(newAlpha)

				c.SetPixel(x, y, newColor)
			}
		}
	}
}

// DrawCircle draws the outline of a circle.
func (c *Canvas) DrawCircle(cx, cy, r int, color uint32) {
	x, y := r, 0
	p := 1 - r

	c.setCirclePixels(cx, cy, x, y, color)

	for x > y {
		y++
		if p <= 0 {
			p += 2*y + 1
		} else {
			x--
			p += 2*(y-x) + 1
		}
		c.setCirclePixels(cx, cy, x, y, color)
	}
}

// TODO: anti aliasing for triangle
// FillTriangle renders a triangle based on given 3 coordinates
func (c *Canvas) FillTriangle(x0, y0, x1, y1, x2, y2 int, color uint32) {
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
			if x >= 0 && x < c.Width && y >= 0 && y < c.Height {
				c.SetPixel(x, y, color)
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
func (c *Canvas) DrawTriangle(x0, y0, x1, y1, x2, y2 int, color uint32) {
	c.Line(x0, y0, x1, y1, color)
	c.Line(x1, y1, x2, y2, color)
	c.Line(x2, y2, x0, y0, color)
}

// TODO: anti aliasing for Line
// Line draws a straight line between 2 points: (x0, y0), (x1, y1)
func (c *Canvas) Line(x0, y0, x1, y1 int, color uint32) {
	dx := abs(x1 - x0)
	sx := -1
	if x0 < x1 {
		sx = 1
	}
	dy := -abs(y1 - y0)
	sy := -1
	if y0 < y1 {
		sy = 1
	}
	err := dx + dy

	for {
		c.SetPixel(x0, y0, color)
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

// PixelsToBytes converts a Canvas's pixel data to a slice of bytes.
func (c *Canvas) PixelsToBytes() []byte {
	if len(c.Pixels) == 0 {
		return nil
	}
	bytes := make([]byte, c.Width*c.Height*4)

	for i, pixel := range c.Pixels {
		j := i * 4
		bytes[j] = byte(pixel >> 24)   // Red
		bytes[j+1] = byte(pixel >> 16) // Green
		bytes[j+2] = byte(pixel >> 8)  // Blue
		bytes[j+3] = byte(pixel)       // Alpha
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

// SaveToPng saves the pixels into the png format, generating a png file
func (c *Canvas) SaveToPng(filePath string) error {
	img := image.NewNRGBA(image.Rect(0, 0, c.Width, c.Height))

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			index := y*c.Stride + x
			pixel := c.Pixels[index]
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

// SaveToPpm saves the pixels into the ppm format, generating a ppm file
func (c *Canvas) SaveToPpm(filePath string) error {
	if filePath == "" {
		return errors.New("file path must be provided")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = fmt.Fprintf(file, "P6\n%d %d\n255\n", c.Width, c.Height)
	if err != nil {
		return err
	}

	for y := 0; y < c.Height; y++ {
		for x := 0; x < c.Width; x++ {
			index := y*c.Stride + x
			pixel := c.Pixels[index]
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

// helper functions

func abs(n int) int {
	if n < 0 {
		return -n
	}

	return n
}

func degreesToRadians(degrees float64) float64 {
	return degrees * (math.Pi / 180)
}

// setCirclePixels sets the pixels for the 8 octants of the circle.
func (c *Canvas) setCirclePixels(cx, cy, x, y int, color uint32) {
	if cx+x < c.Width && cy+y < c.Height {
		c.SetPixel(cx+x, cy+y, color)
	}
	if cx-x >= 0 && cy+y < c.Height {
		c.SetPixel(cx-x, cy+y, color)
	}
	if cx+x < c.Width && cy-y >= 0 {
		c.SetPixel(cx+x, cy-y, color)
	}
	if cx-x >= 0 && cy-y >= 0 {
		c.SetPixel(cx-x, cy-y, color)
	}
	if cx+y < c.Width && cy+x < c.Height {
		c.SetPixel(cx+y, cy+x, color)
	}
	if cx-y >= 0 && cy+x < c.Height {
		c.SetPixel(cx-y, cy+x, color)
	}
	if cx+y < c.Width && cy-x >= 0 {
		c.SetPixel(cx+y, cy-x, color)
	}
	if cx-y >= 0 && cy-x >= 0 {
		c.SetPixel(cx-y, cy-x, color)
	}
}

// TODO: DrawCircle
// TODO: {Draw&&Fill}Elipse
