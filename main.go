package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const CONTENTS_FILENAME = "Contents.json"

func main() {
	start := time.Now()

	contentsFile, err := ioutil.ReadFile(filepath.Join(".", "assets", CONTENTS_FILENAME))
	checkError(err)

	contentsFileContent := ContentFile{}
	err = json.Unmarshal([]byte(contentsFile), &contentsFileContent)
	checkError(err)

	input, err := os.Open("logo.png")
	checkError(err)
	defer input.Close()

	outputDirectory := filepath.Join(".", "output", "AppIcon.appiconset")
	err = os.MkdirAll(outputDirectory, os.ModePerm)
	checkError(err)

	err = ioutil.WriteFile(filepath.Join(outputDirectory, CONTENTS_FILENAME), contentsFile, 0644)
	checkError(err)

	for _, imageItem := range contentsFileContent.Images {
		sizeValueString := splitStringByX(imageItem.Size)[0]
		sizeValue, err := strconv.ParseFloat(sizeValueString, 8)
		checkError(err)

		scaleValueString := splitStringByX(imageItem.Scale)[0]
		scaleValue, err := strconv.ParseFloat(scaleValueString, 8)
		checkError(err)
		log.Println(sizeValue * scaleValue)
	}

	// output, err := os.Create("logo_resized.png")
	// checkError(err)
	// defer output.Close()

	// decodedInput, err := png.Decode(input)
	// checkError(err)

	// inputSpecs := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
	// draw.NearestNeighbor.Scale(inputSpecs, inputSpecs.Rect, decodedInput, decodedInput.Bounds(), draw.Over, nil)
	// png.Encode(output, inputSpecs)

	elapsed := time.Since(start)
	log.Printf("done creating icons in %s", elapsed)
}

func checkError(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func isX(r rune) bool {
	return r == 'x'
}

func splitStringByX(str string) []string {
	return strings.FieldsFunc(str, isX)
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
