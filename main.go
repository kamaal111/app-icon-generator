package main

import (
	_ "embed"
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

//go:embed resources/Contents.json
var contentsFile []byte

func main() {
	start := time.Now()

	outputPath := flag.String("o", "", "output path")
	inputPath := flag.String("i", "", "input path")
	verbose := flag.Bool("v", false, "verbose")

	flag.Parse()

	if *outputPath == "" {
		fmt.Printf("no output path provided\nplease give a output path by giving this command the -o flag with the destination\n")
		os.Exit(1)
	}
	if *inputPath == "" {
		fmt.Printf("no input path provided\nplease give a input path by giving this command the -i flag with the destination\n")
		os.Exit(1)
	}

	contentsFileContent := ContentFile{}
	err := json.Unmarshal(contentsFile, &contentsFileContent)
	checkError(err)

	outputDirectory := filepath.Join(*outputPath, "AppIcon.appiconset")

	os.RemoveAll(outputDirectory)

	err = os.MkdirAll(outputDirectory, os.ModePerm)
	checkError(err)

	err = os.WriteFile(filepath.Join(outputDirectory, "Contents.json"), contentsFile, 0644)
	checkError(err)

	var createdImageNames []string
	var channelsCreated []chan string

	decodedImage, err := decodeImage(*inputPath)
	checkError(err)

	for _, imageItem := range contentsFileContent.Images {
		sizeValueString := splitStringByX(imageItem.Size)[0]
		sizeValue, err := strconv.ParseFloat(sizeValueString, 64)
		checkError(err)

		scaleValueString := splitStringByX(imageItem.Scale)[0]
		scaleValue, err := strconv.ParseFloat(scaleValueString, 64)
		checkError(err)

		scaledSize := sizeValue * scaleValue

		if imageItem.Filename == "" {
			logVerbose(fmt.Sprintf("could not find filename for size of %f", scaledSize), *verbose)
			continue
		}

		if contains(createdImageNames, imageItem.Filename) {
			logVerbose(fmt.Sprintf("file with name %s already created", imageItem.Filename), *verbose)
			continue
		}

		channel := make(chan string)
		channelsCreated = append(channelsCreated, channel)
		go createImage(decodedImage, imageItem, scaledSize, outputDirectory, channel)

		createdImageNames = append(createdImageNames, imageItem.Filename)
	}

	for index, channel := range channelsCreated {
		<-channel
		logVerbose(fmt.Sprintf("created %d out of %d", index+1, len(channelsCreated)), *verbose)
	}

	elapsed := time.Since(start)
	fmt.Printf("done creating icons in %s\n", elapsed)
}

func createImage(decodedImage image.Image, imageItem ImageItem, size float64, outputDirectory string, channel chan string) {
	output, err := os.Create(filepath.Join(outputDirectory, imageItem.Filename))
	checkError(err)
	defer output.Close()

	inputSpecs := image.NewRGBA(image.Rect(0, 0, int(size), int(size)))
	draw.NearestNeighbor.Scale(inputSpecs, inputSpecs.Rect, decodedImage, decodedImage.Bounds(), draw.Over, nil)
	png.Encode(output, inputSpecs)

	channel <- imageItem.Filename
}

func decodeImage(path string) (image.Image, error) {
	input, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	fileExtension := getFileExtension(path)
	switch fileExtension {
	case "jpeg", "jpg":
		return jpeg.Decode(input)
	case "png":
		return png.Decode(input)
	default:
		return nil, fmt.Errorf("%s file extension are not supported", fileExtension)
	}
}

func logVerbose(text string, enabled bool) {
	if enabled {
		fmt.Println(text)
	}
}

func contains(array []string, searchValue string) bool {
	for _, item := range array {
		if item == searchValue {
			return true
		}
	}
	return false
}

func checkError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getFileExtension(str string) string {
	splitString := strings.Split(str, ".")
	fileExtension := splitString[len(splitString)-1]
	return fileExtension
}

func splitStringByX(str string) []string {
	return strings.FieldsFunc(str, func(r rune) bool {
		return r == 'x'
	})
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
