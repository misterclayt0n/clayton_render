package main

import olivego "github.com/misterclayt0n/olive.go/olive"

func main() {
	height := 800
	width := 600

	pixels := olivego.OliveGoFill(height, width, 0xFF00FF00)
	err := olivego.OliveGoSaveToPpm(pixels, "output.ppm")
	if err != nil {
		panic("Failed to save to ppm")
	}
}
