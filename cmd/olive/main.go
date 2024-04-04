package main

import olivego "github.com/misterclayt0n/olive.go/olive"

func main() {
	height := 600
	width := 800
	rectHeight := 200
	rectWidth := 200

	pixels := make([][]uint32, height)
	for i := range pixels {
		pixels[i] = make([]uint32, width)
	}

	olivego.Fill(pixels, width, height, 0xFF202020)
	olivego.FillRect(pixels, width, height, width/2-rectWidth/2, height/2-rectHeight/2, rectWidth, rectHeight, 0xFFFF0000)

	err := olivego.SaveToPpm(pixels, width, height, "output.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}
