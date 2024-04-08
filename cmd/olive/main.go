package main

import olivego "github.com/misterclayt0n/olive.go/olive"

var height int = 800
var width int = 800
var cols int = height / 100
var rows int = width / 100
var cellHeight int = height / cols
var cellWidth int = width / rows
var backgroundColor uint32 = 0xFF202020

// rectHeight := 200
// rectWidth := 200

func buildPixel() [][]uint32 {
	pixels := make([][]uint32, height)
	for i := range pixels {
		pixels[i] = make([]uint32, width)
	}

	return pixels
}

func chessboard() {
	pixels := buildPixel()

	olivego.Fill(pixels, width, height, backgroundColor)

	for y := 0; y < cols; y++ {
		for x := 0; x < rows; x++ {
			var color uint32
			if (x+y)%2 == 0 {
				color = 0xFFFFFFFF
			} else {
				color = 0xFF000000
			}
			olivego.FillRect(pixels, width, height, x*cellWidth, y*cellHeight, cellWidth, cellHeight, color)
		}
	}

	err := olivego.SaveToPpm(pixels, width, height, "chessboard.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func circle() {
	pixels := buildPixel()

	olivego.Fill(pixels, width, height, backgroundColor)
	olivego.FillCircle(pixels, width, height, width/2, height/2, 100, 0xFF00FF00)
	err := olivego.SaveToPpm(pixels, width, height, "circle.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func line() {
	pixels := buildPixel()
	olivego.Fill(pixels, width, height, backgroundColor)
	olivego.Line(pixels, width, height, 0, 0, width, height, 0xFFFF0000)
	olivego.Line(pixels, width, height, width, 0, 0, height, 0xFFFF0000)
	err := olivego.SaveToPpm(pixels, width, height, "line.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func main() {
	circle()
	chessboard()
	line()
}
