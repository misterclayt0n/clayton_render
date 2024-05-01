package main

import (
	"github.com/misterclayt0n/olive.go/olive"
)

const width, height = 800, 600

func main() {
	canvas := olivego.NewCanvas(width, height)
	canvas.Fill(0x181818FF)
	canvas.FillCircle(width/2, height/2, 200, 0xFF000055)
	canvas.SaveToPpm("sexo.ppm")
}
