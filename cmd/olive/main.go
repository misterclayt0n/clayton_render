package main

import olivego "github.com/misterclayt0n/olive.go/olive"

var height int = 800
var width int = 800
var cols int = height / 100
var rows int = width / 100
var cellHeight int = height / cols
var cellWidth int = width / rows
var backgroundColor uint32 = 0x202020FF

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
				color = 0x000000FF
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
	olivego.Line(pixels, width, height, 0, 0, width, height, 0xFF0000FF)
	olivego.Line(pixels, width, height, width, 0, 0, height, 0xFF0000FF)
	err := olivego.SaveToPpm(pixels, width, height, "line.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func triangle() {
	pixels := buildPixel()
	x0, y0 := 200, 200
	x1, y1 := 200, 400
	x2, y2 := 400, 300
	olivego.Fill(pixels, width, height, backgroundColor)
	olivego.FillTriangle(pixels, width, height, x0, y0, x1, y1, x2, y2, 0xFF0000FF)
	err := olivego.SaveToPpm(pixels, width, height, "triangle.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func alpha() {
	pixels := buildPixel()
	olivego.Fill(pixels, width, height, backgroundColor)

	olivego.FillRect(pixels, width, height, 0, 0, width*3/4, height*3/4, 0xFF0000FF)

	for y := height / 4; y < height; y++ {
		for x := width / 4; x < width; x++ {
			if x < width*3/4 && y < height*3/4 {
				olivego.ApplyBlend(pixels, x, y, 0xFF000055)
			}
			olivego.ApplyBlend(pixels, x, y, 0x00FF0055) // 0x00FF00AA é o verde com alpha
		}
	}

	err := olivego.SaveToPpm(pixels, width, height, "alpha.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}

func main() {
	circle()
	chessboard()
	line()
	triangle()
	alpha()
}
