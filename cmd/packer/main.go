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

	"github.com/unitoftime/packer"
)

func main() {
	startTime := time.Now()

	inFlag := flag.String("input", "input", "The directory of the input folder")
	outFlag := flag.String("output", "packed", "The filename of the output json and png")
	extrudeFlag := flag.Int("extrude", 1, "The amount to extrude each sprite")
	statsFlag := flag.Bool("stats", false, "If true, display stats")
	sizeFlag := flag.Int("size", 1024, "The width and height of the packed atlas")
	mountPoint := flag.Bool("mountpoints", false, "If set, the program will analyze color-based mountpoints")
	typeFlag := flag.String("type", "scanline", "Select packer type: scanline, row")
	flag.Parse()

	directory := *inFlag
	output := *outFlag
	extrude := *extrudeFlag
	showStatistics := *statsFlag
	packerType := *typeFlag

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

	if !*mountPoint {
		packer.PrepareImageList(images, extrude)

		// Pack all images
		if packerType == "row" {
			images = packer.RowWisePacker(images, width, height)
		} else {
			images = packer.BasicScanlinePacker(images, width, height)
		}


		imageName := fmt.Sprintf("%s.png", output)

		atlas, data := packer.Pack(imageName, images, width, height)

		jsonFile, err := os.Create(fmt.Sprintf("%s.json", output))
		if err != nil { log.Fatal(err) }

		b, err := json.Marshal(data)
		if err != nil { log.Fatal(err) }
		jsonFile.Write(b)

		outputFile, err := os.Create(imageName)
		if err != nil { log.Fatal(err) }
		png.Encode(outputFile, atlas)
		outputFile.Close()
	} else {
		// Mountpoint mode
		mountFrames := packer.CalculateMountPoints(images)

		jsonFile, err := os.Create(fmt.Sprintf("%s.json", output))
		if err != nil { log.Fatal(err) }

		b, err := json.Marshal(mountFrames)
		if err != nil { log.Fatal(err) }
		jsonFile.Write(b)
	}

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
