package example

import olivego "github.com/misterclayt0n/olive.go/olive"

var (
	height          int    = 800
	width           int    = 800
	cols            int    = height / 100
	rows            int    = width / 100
	cellHeight      int    = height / cols
	cellWidth       int    = width / rows
	backgroundColor uint32 = 0x202020FF
)

// rectHeight := 200
// rectWidth := 200

func Circle() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x55555555)

	canvas.FillCircle(width/2, height/2, 300, 0xFF0000FF)
	if err := canvas.SaveToPpm("circle.ppm"); err != nil {
		panic("shit happens")
	}
}

// func Line() {
// 	pixels := olivego.BuildPixel(height, width)
// 	olivego.Fill(pixels, width, height, backgroundColor)
// 	olivego.Line(pixels, width, height, 0, 0, width, height, 0xFF0000FF)
// 	olivego.Line(pixels, width, height, width, 0, 0, height, 0xFF0000FF)
// 	err := olivego.SaveToPpm(pixels, width, height, "line.ppm")
// 	if err != nil {
// 		panic("Failed to save to ppm")
// 	}
// }

// func Alpha() {
// 	pixels := olivego.BuildPixel(height, width)
// 	olivego.Fill(pixels, width, height, backgroundColor)

// 	olivego.FillRect(pixels, width, height, 0, 0, width*3/4, height*3/4, 0xFF0000FF)
// 	olivego.FillRect(pixels, width, height, width/4, height/4, width*3/4, height*3/4, 0x00FF0055)

// 	err := olivego.SaveToPpm(pixels, width, height, "alpha.ppm")
// 	if err != nil {
// 		panic("Failed to save to ppm")
// 	}
// }

func CanvaPNG() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0xFF000055)
	if err := canvas.SaveToPng("canva.png"); err != nil {
		panic("shit happens")
	}
}

func CangaPPM() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0xFF000055)
	if err := canvas.SaveToPpm("canva.ppm"); err != nil {
		panic("shit happens")
	}
}

func CanvaRect() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x55555555)
	canvas.FillRect(width/2, height/2, width*3/4, height*3/4, 0xFF0000FF)
	if err := canvas.SaveToPpm("canvas_rect.ppm"); err != nil {
		panic("shit happens")
	}
}

func Chessboard() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x55555555)

	for y := 0; y < cols; y++ {
		for x := 0; x < rows; x++ {
			var color uint32
			if (x+y)%2 == 0 {
				color = 0xFFFFFFFF
			} else {
				color = 0x000000FF
			}
			canvas.FillRect(x*cellWidth, y*cellHeight, cellWidth, cellHeight, color)
		}
	}
	if err := canvas.SaveToPpm("chessboard.ppm"); err != nil {
		panic("shit happens")
	}
}

func DrawCircle() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x55555555)
	x0, y0 := 200, 200
	x1, y1 := 200, 400
	x2, y2 := 400, 300

	canvas.DrawRect(x0, y0, width*3/4, height*3/4, 0xFF0000FF)
	canvas.DrawCircle(width/2, height/2, 100, 0xFF0000FF)
	canvas.DrawTriangle(x0, y0, x1, y1, x2, y2, 0xFF0000FF)
	if err := canvas.SaveToPpm("drawing.ppm"); err != nil {
		panic("shit happens")
	}
}

func Triangle() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x55555555)
	x0, y0 := 200, 200
	x1, y1 := 200, 400
	x2, y2 := 400, 300
	canvas.FillTriangle(x0, y0, x1, y1, x2, y2, 0xFF0000FF)
	if err := canvas.SaveToPpm("triangle.ppm"); err != nil {
		panic("shit happens")
	}
}
