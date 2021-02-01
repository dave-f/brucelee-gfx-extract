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
var numberImages []image.RGBA

// Decode a byte into its 2 "left" and "right" pixels
func decodePixel(pixel byte) (l,r byte) {
	l = ((pixel & 0b10) >> 1) | ((pixel & 0b1000) >> 2) | ((pixel & 0b100000) >> 3) | ((pixel & 0b10000000) >> 4)
	r = ((pixel & 0b1) >> 0) | ((pixel & 0b100) >> 1) | ((pixel & 0b10000) >> 2) | ((pixel & 0b1000000) >> 3)
	return
}

// Create the BBC Micro colours
func makeBBCMicroColours() {
	coloursBBC = make([]color.RGBA,16)
	coloursBBC[0] = color.RGBA{0x10, 0x10, 0x10, 0xff} // black
	coloursBBC[1] = color.RGBA{0xff, 0x00, 0x00, 0xff} // red
	coloursBBC[2] = color.RGBA{0x00, 0xff, 0x00, 0xff} // green
	coloursBBC[3] = color.RGBA{0xff, 0xff, 0x00, 0xff} // yellow
	coloursBBC[4] = color.RGBA{0x00, 0x00, 0xff, 0xff} // blue
	coloursBBC[5] = color.RGBA{0xff, 0x00, 0xff, 0xff} // magenta
	coloursBBC[6] = color.RGBA{0x00, 0xff, 0xff, 0xff} // cyan
	coloursBBC[7] = color.RGBA{0xff, 0xff, 0xff, 0xff} // White
	coloursBBC[8] = color.RGBA{0x20, 0x20, 0x20, 0xff} // Last
	coloursBBC[9] = color.RGBA{0x30, 0x30, 0x30, 0xff} // Last
	coloursBBC[10] = color.RGBA{0x40, 0x40, 0x40, 0xff} // Last
	coloursBBC[11] = color.RGBA{0x50, 0x50, 0x50, 0xff} // Last
	coloursBBC[12] = color.RGBA{0x60, 0x60, 0x60, 0xff} // Last
	coloursBBC[13] = color.RGBA{0x70, 0x70, 0x70, 0xff} // Last
	coloursBBC[14] = color.RGBA{0x80, 0x80, 0x80, 0xff} // Last
	coloursBBC[15] = color.RGBA{0x90, 0x90, 0x90, 0xff} // Last
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

func replaceGraphic(b []byte) {
	for i:=0; i<len(b); i++ {
		b[i] &= 0xf
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

// TODO Test getting some numbers out, 0 is at 0,196 and 3,7 bounds
func renderNumberToImage(img *image.RGBA, number int, x int, y int) {

	fmt.Println(strconv.Itoa(number))
	// for each character in this string, render the appropriate character
	// and x+= 4 (each number is 3 pixels)

	//rectForZero := image.Rect(0,196,0+3,196+7)
	//imgZero := img.SubImage(rectForZero)
	//renderFontX := 10
	//renderFontY := 10

	//fmt.Println(rectForZero)
	//fmt.Println("Bounds are ", imgZero.Bounds())

	//for y := 0; y<7; y++ {
	//	for x :=0; x<3; x++ {
	//		img.Set(renderFontX,renderFontY,imgZero.At(x,y+196))
	//		fmt.Println(imgZero.At(x,y))
	//		renderFontX++
	//	}
	//	renderFontX=10
	//	renderFontY++
	//}
}

// TODO Build the font images in `numberImages'
func makeFont(i *image.RGBA, yOffset int) {
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
	pixTableSize := []int{16,16,16,4,4,4,4,4}
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

	ptr := 8006 // lookup table offset

	for i := 0; i < 59; i++ {
		tableAddr := (uint16(data[ptr+1]) << 8) | uint16(data[ptr+0])
		fileOffst := tableAddr + 4096 - 6400
		if (i>0) {
			//outputStr := fmt.Sprintf("Object %02d : Data %04x (File offset %04d), Width %02d bytes (%02d pixels), Height %02d, Pixel Table %02d (%d)", i, tableAddr, fileOffst, data[ptr+2], data[ptr+2]*2, data[ptr+3], data[ptr+4] & 0xfe, data[ptr+4] & 1)
			//fmt.Println(outputStr)
			visualObjectMap[uint32(i)] = VisualObject {
				fileOffst,data[ptr+2],data[ptr+3],(data[ptr+4] & 0xfe)>>1, (data[ptr+4] & 1)==1,
			}
		}
	    ptr += 5
	}

	calcImageHeight := 0

	for _, e := range visualObjectMap {
		calcImageHeight += int(e.heightInRows)
	}

	fmt.Println("Total image height",calcImageHeight)

	// as a test, replace hippo graphic with 0s
	// for i:= 0; i<10*24; i++ {
	//  data[10535+i] &= 0xf
	// }
	// or, do same job:
	// replaceGraphic(data[10535:10535+(10*24)])

	// err = ioutil.WriteFile("new/BRUCE1", data, 0777)

	//if err != nil {
	//	fmt.Println(err)
	//}

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

	for i := 1; i<len(visualObjectMap); i++ {
		// fmt.Println(i,"decoding to image:",visualObjectMap[uint32(i)])
		if (i!=46) { // somthing wrong with this one
			img.Set(24,renderY,grey)
			img.Set(25,renderY,grey)
			img.Set(26,renderY,grey)
			img.Set(27,renderY,grey)
			decodeGraphicToImage(img, visualObjectMap[uint32(i)], 0, renderY)
			renderY += int(visualObjectMap[uint32(i)].heightInRows)
		}
	}

	// TODO Build the font
	makeFont(img,196)

	// TODO Every 5 numbers, draw the count
	renderNumberToImage(img,0,28,0)

	// Save it
	pngFile, _ := os.Create("image.png")
	png.Encode(pngFile, img)

	// printing binary!
	// dave := data[10535]
	// fmt.Println(strconv.FormatInt(int64(dave), 2))
}
