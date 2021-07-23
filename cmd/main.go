package main

import (
	"fmt"
	"log"
	"time"
	"flag"
	"os"
	"io/ioutil"
	"image"
	"image/png"
	"encoding/json"

	"github.com/jstewart7/packer"
)

func main() {
	startTime := time.Now()

	inFlag := flag.String("input", "input", "The directory of the input folder")
	outFlag := flag.String("output", "packed", "The filename of the output json and png")
	extrudeFlag := flag.Int("extrude", 1, "The amount to extrude each sprite")
	statsFlag := flag.Bool("stats", false, "If true, display stats")
	sizeFlag := flag.Int("size", 1024, "The width and height of the packed atlas")
	flag.Parse()

	directory := *inFlag
	output := *outFlag
	extrude := *extrudeFlag
	showStatistics := *statsFlag

	width := *sizeFlag
	height := *sizeFlag

	// Get all images to pack
	images := make([]packer.ImageData, 0)
	files := GetFileList(fmt.Sprintf("./%s/", directory))
	for _, file := range files {
		img, err := LoadImage(fmt.Sprintf("./%s/%s", directory, file))
		if err != nil {
			log.Println("Skipping file (Marked - Not an Image):", file)
			continue
		}
		images = append(images, packer.NewImageData(img, file))
	}

	packer.PrepareImageList(images, extrude)

	// Pack all images
	images = packer.BasicScanlinePacker(images, width, height)

	atlas, data := packer.Pack(images, width, height)

	jsonFile, err := os.Create(fmt.Sprintf("%s.json", output))
	if err != nil { log.Fatal(err) }

	b, err := json.Marshal(data)
	if err != nil { log.Fatal(err) }
	jsonFile.Write(b)

	outputFile, err := os.Create(fmt.Sprintf("%s.png", output))
	if err != nil { log.Fatal(err) }
	png.Encode(outputFile, atlas)
	outputFile.Close()


	if showStatistics {
		packedArea := 0
		for i := range images {
			packedArea += images[i].Area()
		}

		efficiency := float64(packedArea) / float64(width * height)

		fmt.Println("Packing took:", time.Since(startTime))
		fmt.Printf("Efficiency:   %.2f%%\n", 100 * efficiency)
	}
}

func GetFileList(directory string) []string {
	files, err := ioutil.ReadDir(directory)
	if err != nil { panic(err) }

	list := make([]string, 0)
	for _, file := range files {
		list = append(list, file.Name())
	}
	return list
}

func LoadImage(filename string) (image.Image, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	loaded, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}

	return loaded, nil
}
