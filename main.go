package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"golang.org/x/image/draw"
)

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
	err := json.Unmarshal([]byte(CONTENTS), &contentsFileContent)
	checkError(err)

	outputDirectory := filepath.Join(*outputPath, "AppIcon.appiconset")

	os.RemoveAll(outputDirectory)

	err = os.MkdirAll(outputDirectory, os.ModePerm)
	checkError(err)

	err = ioutil.WriteFile(filepath.Join(outputDirectory, "Contents.json"), []byte(CONTENTS), 0644)
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
			logVerbose(fmt.Sprintf("could not find filename for size of %f", scaledSize), *verbose)
			continue
		}

		if contains(createdImageNames, imageItem.Filename) {
			logVerbose(fmt.Sprintf("file with name %s already created", imageItem.Filename), *verbose)
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
		logVerbose(fmt.Sprintf("created %d out of %d", index+1, channelsCreatedLength), *verbose)
	}

	elapsed := time.Since(start)
	fmt.Printf("done creating icons in %s\n", elapsed)
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

const CONTENTS = `
{
	"images" : [
	  {
		"filename" : "40.png",
		"idiom" : "iphone",
		"scale" : "2x",
		"size" : "20x20"
	  },
	  {
		"filename" : "60.png",
		"idiom" : "iphone",
		"scale" : "3x",
		"size" : "20x20"
	  },
	  {
		"filename" : "29.png",
		"idiom" : "iphone",
		"scale" : "1x",
		"size" : "29x29"
	  },
	  {
		"filename" : "58.png",
		"idiom" : "iphone",
		"scale" : "2x",
		"size" : "29x29"
	  },
	  {
		"filename" : "87.png",
		"idiom" : "iphone",
		"scale" : "3x",
		"size" : "29x29"
	  },
	  {
		"filename" : "80.png",
		"idiom" : "iphone",
		"scale" : "2x",
		"size" : "40x40"
	  },
	  {
		"filename" : "120.png",
		"idiom" : "iphone",
		"scale" : "3x",
		"size" : "40x40"
	  },
	  {
		"filename" : "57.png",
		"idiom" : "iphone",
		"scale" : "1x",
		"size" : "57x57"
	  },
	  {
		"filename" : "114.png",
		"idiom" : "iphone",
		"scale" : "2x",
		"size" : "57x57"
	  },
	  {
		"filename" : "120.png",
		"idiom" : "iphone",
		"scale" : "2x",
		"size" : "60x60"
	  },
	  {
		"filename" : "180.png",
		"idiom" : "iphone",
		"scale" : "3x",
		"size" : "60x60"
	  },
	  {
		"filename" : "20.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "20x20"
	  },
	  {
		"filename" : "40.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "20x20"
	  },
	  {
		"filename" : "29.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "29x29"
	  },
	  {
		"filename" : "58.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "29x29"
	  },
	  {
		"filename" : "40.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "40x40"
	  },
	  {
		"filename" : "80.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "40x40"
	  },
	  {
		"filename" : "50.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "50x50"
	  },
	  {
		"filename" : "100.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "50x50"
	  },
	  {
		"filename" : "72.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "72x72"
	  },
	  {
		"filename" : "144.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "72x72"
	  },
	  {
		"filename" : "76.png",
		"idiom" : "ipad",
		"scale" : "1x",
		"size" : "76x76"
	  },
	  {
		"filename" : "152.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "76x76"
	  },
	  {
		"filename" : "167.png",
		"idiom" : "ipad",
		"scale" : "2x",
		"size" : "83.5x83.5"
	  },
	  {
		"filename" : "1024.png",
		"idiom" : "ios-marketing",
		"scale" : "1x",
		"size" : "1024x1024"
	  },
	  {
		"filename" : "48.png",
		"idiom" : "watch",
		"role" : "notificationCenter",
		"scale" : "2x",
		"size" : "24x24",
		"subtype" : "38mm"
	  },
	  {
		"filename" : "55.png",
		"idiom" : "watch",
		"role" : "notificationCenter",
		"scale" : "2x",
		"size" : "27.5x27.5",
		"subtype" : "42mm"
	  },
	  {
		"filename" : "58.png",
		"idiom" : "watch",
		"role" : "companionSettings",
		"scale" : "2x",
		"size" : "29x29"
	  },
	  {
		"filename" : "87.png",
		"idiom" : "watch",
		"role" : "companionSettings",
		"scale" : "3x",
		"size" : "29x29"
	  },
	  {
		"filename" : "66.png",
		"idiom" : "watch",
		"role" : "notificationCenter",
		"scale" : "2x",
		"size" : "33x33",
		"subtype" : "45mm"
	  },
	  {
		"filename" : "80.png",
		"idiom" : "watch",
		"role" : "appLauncher",
		"scale" : "2x",
		"size" : "40x40",
		"subtype" : "38mm"
	  },
	  {
		"filename" : "88.png",
		"idiom" : "watch",
		"role" : "appLauncher",
		"scale" : "2x",
		"size" : "44x44",
		"subtype" : "40mm"
	  },
	  {
		"filename" : "92.png",
		"idiom" : "watch",
		"role" : "appLauncher",
		"scale" : "2x",
		"size" : "46x46",
		"subtype" : "41mm"
	  },
	  {
		"filename" : "100.png",
		"idiom" : "watch",
		"role" : "appLauncher",
		"scale" : "2x",
		"size" : "50x50",
		"subtype" : "44mm"
	  },
	  {
		"filename" : "102.png",
		"idiom" : "watch",
		"role" : "appLauncher",
		"scale" : "2x",
		"size" : "51x51",
		"subtype" : "45mm"
	  },
	  {
		"filename" : "172.png",
		"idiom" : "watch",
		"role" : "quickLook",
		"scale" : "2x",
		"size" : "86x86",
		"subtype" : "38mm"
	  },
	  {
		"filename" : "196.png",
		"idiom" : "watch",
		"role" : "quickLook",
		"scale" : "2x",
		"size" : "98x98",
		"subtype" : "42mm"
	  },
	  {
		"filename" : "216.png",
		"idiom" : "watch",
		"role" : "quickLook",
		"scale" : "2x",
		"size" : "108x108",
		"subtype" : "44mm"
	  },
	  {
		"filename" : "216-1.png",
		"idiom" : "watch",
		"role" : "quickLook",
		"scale" : "2x",
		"size" : "117x117",
		"subtype" : "45mm"
	  },
	  {
		"filename" : "1024.png",
		"idiom" : "watch-marketing",
		"scale" : "1x",
		"size" : "1024x1024"
	  },
	  {
		"filename" : "16.png",
		"idiom" : "mac",
		"scale" : "1x",
		"size" : "16x16"
	  },
	  {
		"filename" : "32.png",
		"idiom" : "mac",
		"scale" : "2x",
		"size" : "16x16"
	  },
	  {
		"filename" : "32.png",
		"idiom" : "mac",
		"scale" : "1x",
		"size" : "32x32"
	  },
	  {
		"filename" : "64.png",
		"idiom" : "mac",
		"scale" : "2x",
		"size" : "32x32"
	  },
	  {
		"filename" : "128.png",
		"idiom" : "mac",
		"scale" : "1x",
		"size" : "128x128"
	  },
	  {
		"filename" : "256.png",
		"idiom" : "mac",
		"scale" : "2x",
		"size" : "128x128"
	  },
	  {
		"filename" : "256.png",
		"idiom" : "mac",
		"scale" : "1x",
		"size" : "256x256"
	  },
	  {
		"filename" : "512.png",
		"idiom" : "mac",
		"scale" : "2x",
		"size" : "256x256"
	  },
	  {
		"filename" : "512.png",
		"idiom" : "mac",
		"scale" : "1x",
		"size" : "512x512"
	  },
	  {
		"filename" : "1024.png",
		"idiom" : "mac",
		"scale" : "2x",
		"size" : "512x512"
	  }
	],
	"info" : {
	  "author" : "xcode",
	  "version" : 1
	}
  }
`
