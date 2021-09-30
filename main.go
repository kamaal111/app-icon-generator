package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

const CONTENTS_FILENAME = "Contents.json"

func main() {
	start := time.Now()

	outputPath := flag.String("o", "", "output path")
	inputPath := flag.String("i", "", "input path")

	flag.Parse()

	if *outputPath == "" {
		fmt.Printf("no output path provided\nplease give a output path by giving this command the -o flag with the destination\n")
		os.Exit(1)
	}
	if *inputPath == "" {
		fmt.Printf("no input path provided\nplease give a input path by giving this command the -i flag with the destination\n")
		os.Exit(1)
	}

	contentsFile, err := ioutil.ReadFile(filepath.Join(".", "assets", CONTENTS_FILENAME))
	checkError(err)

	contentsFileContent := ContentFile{}
	err = json.Unmarshal([]byte(contentsFile), &contentsFileContent)
	checkError(err)

	outputDirectory := filepath.Join(*outputPath, "AppIcon.appiconset")

	os.RemoveAll(outputDirectory)

	err = os.MkdirAll(outputDirectory, os.ModePerm)
	checkError(err)

	err = ioutil.WriteFile(filepath.Join(outputDirectory, CONTENTS_FILENAME), contentsFile, 0644)
	checkError(err)

	var createdImageNames []string
	var channelsCreated []chan string

	for _, imageItem := range contentsFileContent.Images {
		sizeValueString := splitStringByX(imageItem.Size)[0]
		sizeValue, err := strconv.ParseFloat(sizeValueString, 64)
		checkError(err)

		scaleValueString := splitStringByX(imageItem.Scale)[0]
		scaleValue, err := strconv.ParseFloat(scaleValueString, 64)
		checkError(err)

		scaledSize := sizeValue * scaleValue

		if imageItem.Filename == "" {
			log.Printf("could not find filename for size of %f \n", scaledSize)
			continue
		}

		if contains(createdImageNames, imageItem.Filename) {
			log.Printf("file with name %s already created\n", imageItem.Filename)
			continue
		}

		channel := make(chan string)
		channelsCreated = append(channelsCreated, channel)
		go createImage(*inputPath, imageItem, scaledSize, outputDirectory, channel)

		createdImageNames = append(createdImageNames, imageItem.Filename)
	}

	channelsCreatedLength := len(channelsCreated)
	for index, channel := range channelsCreated {
		<-channel
		log.Printf("created %d out of %d\n", index+1, channelsCreatedLength)
	}

	elapsed := time.Since(start)
	log.Printf("done creating icons in %s\n", elapsed)
}

func createImage(inputPath string, imageItem ImageItem, size float64, outputDirectory string, channel chan string) {
	output, err := os.Create(filepath.Join(outputDirectory, imageItem.Filename))
	checkError(err)
	defer output.Close()

	input, err := os.Open(inputPath)
	checkError(err)
	defer input.Close()

	decodedInput, err := png.Decode(input)
	checkError(err)

	inputSpecs := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	draw.NearestNeighbor.Scale(inputSpecs, inputSpecs.Rect, decodedInput, decodedInput.Bounds(), draw.Over, nil)
	png.Encode(output, inputSpecs)

	channel <- imageItem.Filename
}

func contains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
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
