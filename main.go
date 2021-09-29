package main

import (
	"image"
	"image/png"
	"log"
	"os"

	"golang.org/x/image/draw"
)

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

	inputSpecs := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	draw.NearestNeighbor.Scale(inputSpecs, inputSpecs.Rect, decodedInput, decodedInput.Bounds(), draw.Over, nil)
	png.Encode(output, inputSpecs)

	log.Println("done creating icons")
}
