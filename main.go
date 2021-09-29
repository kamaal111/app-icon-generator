package main

import (
	"encoding/json"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"golang.org/x/image/draw"
)

func main() {
	contentsFile, err := ioutil.ReadFile("assets/Contents.json")
	if err != nil {
		log.Fatalln(err)
	}

	contentsFileContent := ContentFile{}
	err = json.Unmarshal([]byte(contentsFile), &contentsFileContent)
	if err != nil {
		log.Fatalln(err)
	}

	for _, imageItem := range contentsFileContent.Images {
		sizeValueString := strings.FieldsFunc(imageItem.Size, func(r rune) bool {
			return r == 'x'
		})[0]
		sizeValue, err := strconv.ParseFloat(sizeValueString, 8)
		if err != nil {
			log.Fatalln(err)
		}
		scaleValueString := strings.FieldsFunc(imageItem.Scale, func(r rune) bool {
			return r == 'x'
		})[0]
		scaleValue, err := strconv.ParseFloat(scaleValueString, 8)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println(sizeValue * scaleValue)
	}

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

type ImageItem struct {
	Filename string `json:"filename"`
	Idiom    string `json:"idiom"`
	Scale    string `json:"scale"`
	Size     string `json:"size"`
}

type ContentInfo struct {
	Author  string `json:"author"`
	Version int    `json:"version"`
}

type ContentFile struct {
	Images []ImageItem `json:"images"`
	Info   ContentInfo `json:"info"`
}
