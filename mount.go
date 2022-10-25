package packer

import (
	"sort"
	"image"
)

//MountPoints[0xFFFFFF]

type MountFrames struct {
	Frames map[string]MountData
}

type MountData struct {
	Filename string
	MountPoints map[uint32]image.Point
}

func CalculateMountPoints(images []ImageData) MountFrames {
	// Sort images by their filename, so that there is some sort of determinism on input of files
	// TODO - this might be automatic, but I haven't tested it cross platform. Maybe check docs
	sort.Slice(images, func(i, j int) bool {
		return images[i].filename < images[j].filename
	})

	frames := MountFrames{
		Frames: make(map[string]MountData),
	}

	//'#00ff00'
	for _, img := range images {
		data := MountData{
			Filename: img.filename,
			MountPoints: make(map[uint32]image.Point),
		}

		bounds := img.img.Bounds()
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
				color := img.img.At(x, y)
				r, g, b, a := color.RGBA()
				if a == 0 {continue } // Skip if the pixel is transparent

				// fmt.Println("Found one:")
				// fmt.Printf("#%x%x%x\n", r, g, b)
				// fmt.Println(r, g, b, a)
				// fmt.Println(color)
				colorKey := ((r>>8) << 16) | ((g>>8) << 8) | ((b>>8) << 4)
				// fmt.Printf("%x\n", colKey)

				data.MountPoints[colorKey] = image.Point{
					(bounds.Dx()/2) - x,
					(bounds.Dy()/2) - y,
				}
			}
			frames.Frames[data.Filename] = data
		}
	}

	return frames
}
