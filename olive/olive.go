package olivego

import (
	"errors"
	"fmt"
	"os"
)

func Fill(pixels [][]uint32, width, height int, color uint32) {
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			pixels[y][x] = color
		}
	}
}

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
			_, err := file.Write([]byte{byte(pixel >> 16), byte(pixel >> 8), byte(pixel)})
			if err != nil {
				panic("could not write to ppm")
			}
		}
	}

	return nil
}

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

func Lines(pixels [][]uint32, width, height int, x1, x2, y1, y2 int, color uint32) {
	// equação da reta: y = mx + n
	dx := x2 - x1
	dy := y2 - y1
	if dx != 0 {
		m := dy / dx
	} else {
		// TODO: not implemented
	}
}
