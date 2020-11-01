# Packer
This is a work in progress texture packing standalone. Right now the packing algorithm isn't very optimized in terms of speed, but probably finds a fairly well packed set of images.

### Install
```
go get github.com/jstewart7/packer
```
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

### Remaining Work
* Create package to let developers dynamically pack images (ie not through command line)
* Dynamically resize atlas image (rather than using a size flag)
* Optimize
