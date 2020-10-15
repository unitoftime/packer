package main

import (
	"log"
	"fmt"
	"os"
	"image"
	"image/draw"
	"image/png"
	"sort"

	"io/ioutil"
	"encoding/json"
	"flag"
)

func main() {
	inFlag := flag.String("input", "input", "The directory of the input folder")
	outFlag := flag.String("output", "packed", "The filename of the output json and png")
	extrudeFlag := flag.Int("extrude", 1, "The amount to extrude each sprite")
	flag.Parse()

	directory := *inFlag
	output := *outFlag
	extrude := *extrudeFlag

	width := 1024
	height := 1024

	// Get all images to pack
	images := make([]ImageData, 0)
	files := GetFileList(fmt.Sprintf("./%s/", directory))
	for _, file := range files {
		img := LoadImage(fmt.Sprintf("./%s/%s", directory, file))
		images = append(images, NewImageData(img, file))
	}

	for i := range images {
		images[i].img = ExtrudeImage(images[i].img, extrude)
		// images[i].extrudeBounds = images[i].img.Bounds()
		images[i].extrudeOffset = image.Point{extrude, extrude}
	}

	// Pack all images
	images = NaiveGreedyPacker(images, width, height)

	atlas, data := Pack(images, width, height)

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

type ImageData struct {
	img image.Image
	filename string
	position image.Point
	origBounds image.Rectangle
	// extrudeBounds image.Rectangle
	extrudeOffset image.Point
}

func NewImageData(img image.Image, filename string) ImageData {
	return ImageData{
		img: img,
		filename: filename,
		position: image.Point{},
		origBounds: img.Bounds(),
	}
}

func NaiveGreedyPacker(images []ImageData, width, height int) []ImageData {
	// 1. Sort rectangles based on order
	// 2. loop through width,height rectangle and place them at first available position

	// Sort
	// TODO - test
	sort.Slice(images, func(i, j int) bool { return (images[i].img.Bounds().Size().X * images[i].img.Bounds().Size().Y) < (images[j].img.Bounds().Size().X * images[j].img.Bounds().Size().Y) })

	targetBounds := image.Rect(0, 0, width, height) // placed image must fall inside the targetBounds

	placed := make([]ImageData, 0)
	// Place Greedily
	for i := range images {

	attempt:
		for y := 0; y < height; y++ {
			for x := 0; x < width; x++ {

				// Check if we can place it here
				attemptPos := image.Point{x,y}
				attemptRect := images[i].img.Bounds().Add(attemptPos)
				if !attemptRect.In(targetBounds) {
					// If the attempt rectangle isn't fully inside the target rect, then fail this position
					continue
				}

				success := true
				for _,placedImg := range placed {
					placedRect := placedImg.img.Bounds().Add(placedImg.position)
					if attemptRect.Overlaps(placedRect) {
						// If there is ever an overlap then break

						// However, we can safely increment X to the point after the image
						x = placedRect.Max.X

						success = false
						break
					}
				}

				if !success { continue }

				// If we were successful in placing, then place it officially
				images[i].position = attemptPos
				placed = append(placed, images[i])
				break attempt
			}
		}
	}

	return placed
}

func Pack(images []ImageData, width, height int) (image.Image, SerializedSpritesheet) {
	data := SerializedSpritesheet{
		Frames: make(map[string]SerializedFrame),
		Meta: make(map[string]interface{}),
	}
	data.Meta["protocol"] = "github.com/jstewart7/packer"

	atlasBounds := image.Rect(0, 0, width, height)
	atlas := image.NewNRGBA(atlasBounds)

	currentBounds := image.Rectangle{}
	// currentPos := image.Point{}
	for _, imageData := range images {
		// img := imageData.img
		// origBounds := img.Bounds()

		// Extrude image
		// img = ExtrudeImage(img, extrude)
		// extrudeBounds := img.Bounds()

		// destOrigBounds := origBounds.Add(currentPos).Add(image.Point{extrude,extrude})
		// destBounds := extrudeBounds.Add(currentPos)
		// draw.Draw(atlas, destBounds, img, image.ZP, draw.Src)
		// currentPos.X += extrudeBounds.Dx()

		img := imageData.img
		extrudeBounds := img.Bounds()
		destOrigBounds := imageData.origBounds.Add(imageData.position).Add(imageData.extrudeOffset)
		destBounds := extrudeBounds.Add(imageData.position)
		draw.Draw(atlas, destBounds, img, image.ZP, draw.Src)
		// currentPos.X += extrudeBounds.Dx()


		currentBounds = currentBounds.Union(destBounds)

		data.Frames[imageData.filename] = SerializedFrame{
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

	// TODO - shrink final atlas down if possible

	return atlas, data
}

// TODO - this is inefficient, but might not matter that much. I think most people will only extrude once
func ExtrudeImage(img image.Image, extrude int) image.Image {
	for i := 0; i < extrude; i++ {
		img = ExtrudeImageOnce(img)
	}
	return img
}

// TODO - needs cleanup
func ExtrudeImageOnce(img image.Image) image.Image {
	extrude := 1
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
