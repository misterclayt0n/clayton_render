package example

import (
	"unsafe"

	"github.com/misterclayt0n/olive.go/olive"
	"github.com/veandco/go-sdl2/sdl"
)

const sdlWidth, sdlHeight int32 = 800, 600

func OliveSDL() {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("OliveGo SDL Window", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, sdlWidth, sdlHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	texture, err := renderer.CreateTexture(uint32(sdl.PIXELFORMAT_RGBA32), sdl.TEXTUREACCESS_STREAMING, sdlWidth, sdlHeight)
	if err != nil {
		panic(err)
	}
	defer texture.Destroy()

	texture.SetBlendMode(sdl.BLENDMODE_BLEND)

	angle := 0.0
	cx, cy := float64(sdlWidth)/2, float64(sdlHeight)/2
	circleX, circleY := cx, cy
	circleVelX, circleVelY := 5.0, 5.0
	radius := 100

	running := true
	for running {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				running = false
			}
		}

		pixels := olivego.BuildPixel(int(sdlHeight), int(sdlWidth))

		x1, y1 := float64(sdlWidth)/2, float64(sdlHeight)/4
		x2, y2 := float64(sdlWidth)/4, float64(3*sdlHeight)/4
		x3, y3 := float64(3*sdlWidth)/4, float64(3*sdlHeight)/4

		newX1, newY1 := olivego.RotatePoint(cx, cy, x1, y1, angle)
		newX2, newY2 := olivego.RotatePoint(cx, cy, x2, y2, angle)
		newX3, newY3 := olivego.RotatePoint(cx, cy, x3, y3, angle)

		olivego.Fill(pixels, int(sdlWidth), int(sdlHeight), 0x202020FF)

		olivego.FillTriangle(pixels, int(sdlWidth), int(sdlHeight),
			int(newX1), int(newY1), int(newX2), int(newY2), int(newX3), int(newY3), 0xFF0000FF)

		circleX += circleVelX
		circleY += circleVelY

		if circleX-float64(radius) <= 0 || circleX+float64(radius) >= float64(sdlWidth) {
			circleVelX = -circleVelX
		}

		if circleY-float64(radius) <= 0 || circleY+float64(radius) >= float64(sdlHeight) {
			circleVelY = -circleVelY
		}

		olivego.FillCircle(pixels, width, height, int(circleX), int(circleY), radius, 0x00FF0055)
		// olivego.FillRect(pixels, int(sdlWidth), int(sdlHeight), int(sdlWidth)/2, int(sdlHeight)/2, int(sdlWidth)/2, int(sdlHeight)/2, 0x0000FF55)
		// olivego.DrawCircle(pixels, width, height, int(sdlWidth)/2, int(sdlHeight)/2, 100, 0x00FF00FF)

		angle += 5.0

		bytePixels := olivego.PixelsToBytes(pixels)

		if len(bytePixels) == 0 {
			bytePixels = make([]byte, sdlWidth*sdlHeight*4)
		}

		texture.Update(nil, unsafe.Pointer(&bytePixels[0]), int(sdlWidth)*4)

		renderer.Clear()
		renderer.Copy(texture, nil, nil)
		renderer.Present()

		sdl.Delay(16)
	}
}
