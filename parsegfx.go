// Bruce Lee BBC Micro graphics extractor
// TODO Explain how all graphics are all stored

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"image"
	"image/png"
	"image/color"
	"strconv"
)

const TotalGraphicsObjects = 58
const FontWidth = 3
const FontHeight = 7

// The entire file
var data []byte

// A map of all the visual objects
type VisualObject struct {
	fileOffset uint16
	widthInBytes byte
	heightInRows byte
	pixelTableIndex byte
	maskFlag bool
}

var visualObjectMap map[uint32]VisualObject

// There are 8 pixel lookup tables
var pixelTables [][]byte

// Define the BBC Micro colours, so we can put them into the image
var coloursBBC []color.RGBA

// Define a simple 0-9 character set
var numberImages [10][3][7]color.RGBA

// Decode a byte into its 2 "left" and "right" pixels
func decodePixel(pixel byte) (l,r byte) {
	l = ((pixel & 0b10) >> 1) | ((pixel & 0b1000) >> 2) | ((pixel & 0b100000) >> 3) | ((pixel & 0b10000000) >> 4)
	r = ((pixel & 0b1) >> 0) | ((pixel & 0b100) >> 1) | ((pixel & 0b10000) >> 2) | ((pixel & 0b1000000) >> 3)
	return
}

// Create the BBC Micro colours
func makeBBCMicroColours() {
	coloursBBC = make([]color.RGBA,16)
	coloursBBC[0] = color.RGBA{0x00, 0x00, 0x00, 0xff} // black
	coloursBBC[1] = color.RGBA{0xff, 0x00, 0x00, 0xff} // red
	coloursBBC[2] = color.RGBA{0x00, 0xff, 0x00, 0xff} // green
	coloursBBC[3] = color.RGBA{0xff, 0xff, 0x00, 0xff} // yellow
	coloursBBC[4] = color.RGBA{0x00, 0x00, 0xff, 0xff} // blue
	coloursBBC[5] = color.RGBA{0xff, 0x00, 0xff, 0xff} // magenta
	coloursBBC[6] = color.RGBA{0x00, 0xff, 0xff, 0xff} // cyan
	coloursBBC[7] = color.RGBA{0xff, 0xff, 0xff, 0xff} // white
	coloursBBC[8] = color.RGBA{0x20, 0x20, 0x20, 0xff} // black 1
	coloursBBC[9] = color.RGBA{0x7f, 0x00, 0x00, 0xff} // red 1
	coloursBBC[10] = color.RGBA{0x00, 0x7f, 0x00, 0xff} // green 1
	coloursBBC[11] = color.RGBA{0x7f, 0x7f, 0x00, 0xff} // yellow 1
	coloursBBC[12] = color.RGBA{0x00, 0x00, 0x7f, 0xff} // blue 1
	coloursBBC[13] = color.RGBA{0x7f, 0x00, 0x7f, 0xff} // magenta 1
	coloursBBC[14] = color.RGBA{0x00, 0x7f, 0x7f, 0xff} // cyan 1
	coloursBBC[15] = color.RGBA{0x7f, 0x7f, 0x7f, 0xff} // white 1
}

func printGraphicsObject(name string, b []byte, pixelTable []byte, w int, h int, flag bool) {
	fmt.Println(name)
	for i:= 0; i<len(b); i++ {
		thisByte := b[i]
		if (flag) {
			thisByte >>= 4
		}
		thisByte &= 0xf
		l,r := decodePixel(pixelTable[thisByte])
		outputStr := fmt.Sprintf("%02x -> %02x / %02x", thisByte, l, r)
		fmt.Println(outputStr)
	}
}

func printPixelTable(b []byte) {
	for i := 0; i<len(b); i++ {
		l,r := decodePixel(b[i])
		outputStr := fmt.Sprintf("%02x -> %02x / %02x", i, l, r)
		fmt.Println(outputStr)
	}
}

func replaceGraphic(b []byte, upperNibble bool) {
	for i:=0; i<len(b); i++ {
		originalByte := b[i]
		if (upperNibble) {
			originalByte &= 0x0f
			originalByte |= 1<<4
			b[i] = originalByte
		} else {
			originalByte &= 0xf
			originalByte |= 1
			b[i] = originalByte
		}
	}
}

func decodeGraphicToImage(i* image.RGBA, o VisualObject, x int, y int) {
	var h byte
	var w byte
	offset := o.fileOffset
	ourx := x
	oury := y

	for w = 0; w < o.widthInBytes; w++ {
		for h = 0; h < o.heightInRows; h++ {
			dataByte := data[offset]
			if (o.maskFlag) {
				dataByte = dataByte >> 4
			}

			dataByte &= 0xf

			actualPixelByte := pixelTables[o.pixelTableIndex][dataByte]
			l,r := decodePixel(actualPixelByte)

			i.Set(ourx, oury, coloursBBC[l])
			i.Set(ourx+1, oury, coloursBBC[r])

			offset++
			oury++
		}
		oury = y
		ourx += 2
	}
}

func renderNumberToImage(img *image.RGBA, number int, x int, y int) {

	numString := strconv.Itoa(number)

	curX := x
	curY := y

	for i:=0;i<len(numString);i++ {
		thisNumber,_ := strconv.Atoi(numString[i:i+1])
		for y:=0;y<FontHeight;y++ {
			for x:=0;x<FontWidth;x++ {
				img.Set(curX,curY,numberImages[thisNumber][x][y])
				curX++
			}
			curX=x
			curY++
		}
		x+=4
		curX=x
		curY=y
	}
}

func makeFont(srcImg *image.RGBA, yOffset int) {

	curY := yOffset

	for i:=0;i<10;i++ {
		for y:=0;y<FontHeight;y++ {
			for x:=0;x<FontWidth;x++ {
				numberImages[i][x][y] = srcImg.RGBAAt(x,curY)
			}
			curY++
		}
	}
}

func main() {

	if len(os.Args) != 3 {
		//fmt.Println("Usage: parsegfx <filename>")
		//return
	}

	visualObjectMap = make(map[uint32]VisualObject)

	f, err := os.Open("org/BRUCE1")//os.Args[1])

	if err != nil {
		fmt.Println(err)
		return
	}

	defer f.Close()
	data, err = ioutil.ReadAll(f)

	if err != nil {
		fmt.Println(err)
		return
	}

	// Build the 8 pixel lookup tables
	pixTableOffs := 7620
	pixTableSize := []int{16,16,16,4,4,4,4,16}
	pixelTables = make([][]byte,8)

	for i, e := range pixTableSize {
		//fmt.Println("Pixel table",i+1)
		//printPixelTable(data[pixTableOffs:pixTableOffs+e])
		//fmt.Println()
		pixelTables[i] = data[pixTableOffs:pixTableOffs+e]
		pixTableOffs += e
	}

	makeBBCMicroColours()

	//sanity check
	//for i, e := range pixelTables {
	//	fmt.Println(i,e)
	//}

	ptr := 8006 + 5 // lookup table offset, the first five bytes are all 0

	for i := 0; i < TotalGraphicsObjects; i++ {
		tableAddr := (uint16(data[ptr+1]) << 8) | uint16(data[ptr+0])
		fileOffst := tableAddr + 4096 - 6400
		//outputStr := fmt.Sprintf("Object %02d : Data %04x (File offset %04d), Width %02d bytes (%02d pixels), Height %02d, Pixel Table %02d (%d)", i, tableAddr, fileOffst, data[ptr+2], data[ptr+2]*2, data[ptr+3], data[ptr+4] & 0xfe, data[ptr+4] & 1)
		//fmt.Println(outputStr)
		visualObjectMap[uint32(i)] = VisualObject {
			fileOffst,data[ptr+2],data[ptr+3],(data[ptr+4] & 0xfe)>>1, (data[ptr+4] & 1)==1,
		}
	    ptr += 5
	}

	fmt.Println("Total objects", len(visualObjectMap))

	calcImageHeight := 0

	for _, e := range visualObjectMap {
		calcImageHeight += int(e.heightInRows)
	}

	fmt.Println("Total image height",calcImageHeight)

	// Create an image in which to display our parsed graphics
	const imageWidth = 640
	imageHeight := calcImageHeight
	black := color.RGBA{0, 0, 0, 0xff}
	grey := color.RGBA{0x7f,0x7f,0x7f,0xff}
	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth,imageHeight}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Clear it
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			img.Set(x, y, black)
		}
	}

	// Render the graphics in, and draw a little indicator at each one
	renderY := 0

	for i := 0; i<len(visualObjectMap); i++ {
		img.Set(24,renderY,grey)
		img.Set(25,renderY,grey)
		img.Set(26,renderY,grey)
		img.Set(27,renderY,grey)
		decodeGraphicToImage(img, visualObjectMap[uint32(i)], 0, renderY)
		renderY += int(visualObjectMap[uint32(i)].heightInRows)
	}

	// This lets us build our 0-9 characters..
	makeFont(img,196)

	// ..so we can draw an ID every 5 objects
	renderY = 0
	for i := 0; i<len(visualObjectMap); i++ {
		if (i%5==0) {
			renderNumberToImage(img,i,30,renderY)
		}
		fmt.Println(i, ":", visualObjectMap[uint32(i)])
		renderY += int(visualObjectMap[uint32(i)].heightInRows)
	}

	// Save it
	pngFile, _ := os.Create("image.png")
	png.Encode(pngFile, img)
	defer pngFile.Close()

	// For object 34, and probably more, the first one is used as a "disappearing" tile, in that the physical colour
	// can be changed from 9 to 0, so the tile can easily be made to disappear
	// replaceGraphic(data[visualObjectMap[34].fileOffset:visualObjectMap[34].fileOffset + uint16(visualObjectMap[34].widthInBytes * visualObjectMap[34].heightInRows)], false)
	err = ioutil.WriteFile("new/BRUCE1", data, 0777)
	if err != nil {
		fmt.Println(err)
	}

	// printing binary!
	// dave := data[10535]
	// fmt.Println(strconv.FormatInt(int64(dave), 2))
}
