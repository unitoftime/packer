package main

import (
	"log"
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/png"
	"io/ioutil"
	"encoding/json"
	"flag"
)

func main() {
	inFlag := flag.String("input", "originals", "The directory of the input folder")
	outFlag := flag.String("output", "packed", "The filename of the output json and png")
	extrudeFlag := flag.Int("extrude", 1, "The amount to extrude each sprite")
	flag.Parse()

	directory := *inFlag
	output := *outFlag
	extrude := *extrudeFlag

	size := 1024
	atlasBounds := image.Rect(0, 0, size, 32)
	atlas := image.NewNRGBA(atlasBounds)

	data := SerializedSpritesheet{
		Frames: make(map[string]SerializedFrame),
		Meta: make(map[string]interface{}),
	}
	data.Meta["protocol"] = "github.com/jstewart7/packer"

	currentBounds := image.Rectangle{}

	currentPos := image.Point{0,0}
	files := GetFileList(fmt.Sprintf("./%s/", directory))
	for _, file := range files {
		img := LoadImage(fmt.Sprintf("./%s/%s", directory, file))
		origBounds := img.Bounds()
		img = ExtrudeImage(img, extrude)
		extrudeBounds := img.Bounds()

		destOrigBounds := origBounds.Add(currentPos).Add(image.Point{extrude,extrude})
		destBounds := extrudeBounds.Add(currentPos)
		draw.Draw(atlas, destBounds, img, image.ZP, draw.Src)
		currentPos.X += extrudeBounds.Dx()

		currentBounds = currentBounds.Union(destBounds)

		data.Frames[file] = SerializedFrame{
			Frame: SerializedRect{
				X: float64(destOrigBounds.Min.X),
				Y: float64(destOrigBounds.Min.Y),
				W: float64(destOrigBounds.Dx()),
				H: float64(destOrigBounds.Dy()),
			},
			Rotated: false, // TODO
			Trimmed: false, // TODO
			SpriteSourceSize: SerializedRect{
				// TODO
			},
			SourceSize: SerializedDim{
				// TODO
			},
			Pivot: SerializedPos{
				// TODO
			},
		}
	}
	log.Println(currentBounds)

	jsonFile, err := os.Create(fmt.Sprintf("%s.json", output))
	if err != nil { log.Fatal(err) }

	b, err := json.Marshal(data)
	if err != nil { log.Fatal(err) }
	jsonFile.Write(b)

	outputFile, err := os.Create(fmt.Sprintf("%s.png", output))
	if err != nil { log.Fatal(err) }
	png.Encode(outputFile, atlas)
	outputFile.Close()
}

// TODO - only extrudes once
func ExtrudeImage(img image.Image, extrude int) image.Image {
	bounds := img.Bounds()
	newImg := image.NewNRGBA(image.Rect(0, 0, bounds.Dx() + (2 * extrude), bounds.Dy() + (2 * extrude)))
	dstBounds := newImg.Bounds()

	draw.Draw(newImg, bounds.Add(image.Point{extrude,extrude}), img, image.ZP, draw.Src)

	// Outer Rows
	ySrc := 0
	yDst := 0
	for xSrc := 0; xSrc < bounds.Dx(); xSrc++ {
		xDst := xSrc+1
		newImg.Set(xDst, yDst, img.At(xSrc, ySrc))
	}

	ySrc = bounds.Dy()-1
	yDst = dstBounds.Dy()-1
	for xSrc := 0; xSrc < bounds.Dx(); xSrc++ {
		xDst := xSrc+1
		newImg.Set(xDst, yDst, img.At(xSrc, ySrc))
	}

	// Corners
	newImg.Set(dstBounds.Min.X, dstBounds.Min.Y, img.At(bounds.Min.X, bounds.Min.Y))
	newImg.Set(dstBounds.Max.X-1, dstBounds.Min.Y, img.At(bounds.Max.X-1, bounds.Min.Y))
	newImg.Set(dstBounds.Min.X, dstBounds.Max.Y-1, img.At(bounds.Min.X, bounds.Max.Y-1))
	newImg.Set(dstBounds.Max.X-1, dstBounds.Max.Y-1, img.At(bounds.Max.X-1, bounds.Max.Y-1))

	// Outer Columns
	xSrc := 0
	xDst := 0
	for ySrc := 0; ySrc < bounds.Dy(); ySrc++ {
		yDst := ySrc+1
		newImg.Set(xDst, yDst, img.At(xSrc, ySrc))
	}

	xSrc = bounds.Dx()-1
	xDst = dstBounds.Dx()-1
	for ySrc := 0; ySrc < bounds.Dx(); ySrc++ {
		yDst := ySrc+1
		newImg.Set(xDst, yDst, img.At(xSrc, ySrc))
	}

	return newImg
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

func LoadImage(filename string) image.Image {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal("Error Opening File: ", filename, " ", err)
	}
	defer file.Close()

	loaded, _, err := image.Decode(file)
	if err != nil {
		log.Fatal("Error Decoding File: ", filename, " ",  err)
	}

	return loaded
}

type SerializedRect struct {
	X,Y,W,H float64
}
type SerializedPos struct {
	X,Y float64
}
type SerializedDim struct {
	W,H float64
}

type SerializedFrame struct {
	Frame SerializedRect
	Rotated bool
	Trimmed bool
	SpriteSourceSize SerializedRect
	SourceSize SerializedDim
	Pivot SerializedPos
}
type SerializedSpritesheet struct {
	Frames map[string]SerializedFrame
	Meta map[string]interface{}
}
