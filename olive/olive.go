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
