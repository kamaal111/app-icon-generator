package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/draw"
)

const LOGO_SIZE = 800

func main() {
	input, err := os.Open("logo.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer input.Close()

	output, err := os.Create("logo_resized.png")
	if err != nil {
		log.Fatalln(err)
	}
	defer output.Close()

	decodedInput, err := png.Decode(input)
	if err != nil {
		log.Fatalln(err)
	}

	log.Println(decodedInput.Bounds().Max.X)
	log.Println(decodedInput.Bounds().Max.Y)

	inputSpecs := image.NewRGBA(image.Rect(0, 0, decodedInput.Bounds().Max.X/2, decodedInput.Bounds().Max.Y/2))
	draw.NearestNeighbor.Scale(inputSpecs, inputSpecs.Rect, decodedInput, decodedInput.Bounds(), draw.Over, nil)
	png.Encode(output, inputSpecs)
}
