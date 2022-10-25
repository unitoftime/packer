# Packer
This is a work in progress texture packing standalone. Right now the packing algorithm isn't very optimized in terms of speed, but probably finds a fairly well packed set of images.

### Install
```
go install github.com/unitoftime/packer/cmd/packer@latest
```
Note: the binary is held in `cmd/packer/` so the suffix `/...` is required to install the `packer` binary to your go path.
### Usage
Basic
```
packer --input sprites --output ./path/to/file
```
#### Flags
```
--input <Directory> - The directory of the input folder
--output <Filename> - The filename of the output json and png files
--extrude <Value> - The amount to extrude each sprite
--stats - If true, display statistics
--size <Value> - The width and height of the packed atlas
```

#### Serialized Json
The JSON file created by the `packer` binary is of form `SerializedSpritesheet`. (See `packer.go` for reference):
```
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
	ImageName string
	Frames map[string]SerializedFrame
	Meta map[string]interface{}
}
```

### Example Output
You can go to the `packer/cmd/images/` and use the `generate.sh` script to generate some random images. Then you can rerun the test.

![packed](https://user-images.githubusercontent.com/2606873/126796465-2203321e-729f-4811-85e3-8b8ee0661d4e.png)


### Remaining Work
* Create package to let developers dynamically pack images (ie not through command line)
* Dynamically resize atlas image (rather than using a size flag)
* Optimize
