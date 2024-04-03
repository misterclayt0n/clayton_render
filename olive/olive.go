package olivego

import (
	"errors"
	"fmt"
	"os"
)

func OliveGoFill(height, width int, color uint32) [][]uint32 {
	pixels := make([][]uint32, height)

	for y := range pixels {
		pixels[y] = make([]uint32, width)
		for x := range pixels[y] {
			pixels[y][x] = color
		}
	}

	return pixels
}

func OliveGoSaveToPpm(pixels [][]uint32, filePath string) error {
	if filePath == "" {
		return errors.New("File path must be provided")
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	height := len(pixels)
	width := 0

	if height > 0 {
		width = len(pixels[0])
	}

	_, err = fmt.Fprintf(file, "P6\n%d %d\n255\n", width, height)
	if err != nil {
		return err
	}

	for _, row := range pixels {
		for _, col := range row {
			_, err := file.Write([]byte{byte(col >> 16), byte(col >> 8), byte(col)})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
