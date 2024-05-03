package main

import (
	"github.com/misterclayt0n/clayton_render/clayton_render"
)

const width, height = 800, 600

func main() {
	canvas := clayton_render.NewCanvas(width, height)
	canvas.Fill(0x181818FF)
	canvas.FillCircle(width/2, height/2, 200, 0xFF000055)
	canvas.SaveToPpm("sexo.ppm")
}
