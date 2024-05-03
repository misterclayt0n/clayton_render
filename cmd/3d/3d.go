package main

import (
	"unsafe"

	"github.com/misterclayt0n/clayton_render/clayton_render"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	width, height int32  = 800, 600
	gridCount            = 10
	gridPad              = 1. / float64(gridCount)
	gridSize             = ((gridCount - 1) * gridPad)
	circleRadius  int    = 10
	circleColor   uint32 = 0xFF2020FF
)

func main() {
	render()
}

func render3d() *clayton_render.Canvas {
	canvas := clayton_render.NewCanvas(int(width), int(height))
	canvas.Fill(0x202020FF)

	fov := 1000.0
	viewerDistance := 5.0

	for cx := 0; cx < gridCount; cx++ {
		for cy := 0; cy < gridCount; cy++ {
			for cz := 0; cz < 1; cz++ {
				x := float64(cx)*gridPad - gridSize/2
				y := float64(cy)*gridPad - gridSize/2
				z := float64(cz)*gridPad - gridSize/2

				// Aplicando a projeção perspectiva.
				scaleFactor := fov / (fov + z*viewerDistance)
				xProjected := x * scaleFactor
				yProjected := y * scaleFactor

				// Convertendo coordenadas projetadas para posições no canvas.
				xCanvas := int((xProjected + 1) / 2 * float64(width))
				yCanvas := int((yProjected + 1) / 2 * float64(height))

				// Ajustando o raio do círculo com base no fator de escala para simular profundidade.
				radius := int(float64(circleRadius) * scaleFactor)

				canvas.FillCircle(xCanvas, yCanvas, radius, circleColor)
			}
		}
	}

	return canvas
}

func render() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("clayton_render 3D", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, width, height, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, width, height)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	// canvas := clayton_render.NewCanvas(int(width), int(height))

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		canvas := render3d()

		bytePixels := canvas.PixelsToBytes()
		if len(bytePixels) == 0 {

		}
		texture.Update(nil, unsafe.Pointer(&bytePixels[0]), int(width)*4)

		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()
		sdl.Delay(16)
	}
}
