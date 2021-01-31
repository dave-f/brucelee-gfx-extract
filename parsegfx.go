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
	//"strconv"
)

// The game contains 8 pixel lookup tables
// 

// A map of all the visual objects
type VisualObject struct {
	fileOffset uint16
	widthInBytes byte
	heightInRows byte
	pixelTableIndex byte
	maskFlag bool
}

var visualObjectMap map[uint32]VisualObject
var pixelTables [][]byte

func decodePixel(pixel byte) (l,r byte) {
	l = ((pixel & 0b10) >> 1) | ((pixel & 0b1000) >> 2) | ((pixel & 0b100000) >> 3) | ((pixel & 0b10000000) >> 4)
	r = ((pixel & 0b1) >> 0) | ((pixel & 0b100) >> 1) | ((pixel & 0b10000) >> 2) | ((pixel & 0b1000000) >> 3)
	return
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

func decodeGraphicToImage(o VisualObject) {
	fmt.Println(o)
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
	data, err := ioutil.ReadAll(f)

	if err != nil {
		fmt.Println(err)
		return
	}

	//var totalBytes = len(data)
	//fmt.Println(totalBytes, "total bytes")

	// the graphics objects are in y order
	// printGraphicsObject("brick",data[10535:10535+(3*10)],data[7620:7620+16],3,10,false)
	// printGraphicsObject("hippo",data[10535:10535+(10*24)],data[7620:7620+16],10,24,true)
	// printGraphicsObject("green statue",16)
	// printGraphicsObject("lantern",17)

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
				fileOffst,data[ptr+2],data[ptr+3],data[ptr+4] & 0xfe, (data[ptr+4] & 1)==1,
			}
		}
	    ptr += 5
	}

	//for i, e := range visualObjectMap {
	//	fmt.Println(i,e)
	//}

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
	const imageHeight = 480
	black := color.RGBA{0, 0, 0, 0xff}
	upLeft := image.Point{0, 0}
	lowRight := image.Point{imageWidth,imageHeight}
	img := image.NewRGBA(image.Rectangle{upLeft, lowRight})

	// Clear it
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			img.Set(x, y, black)
		}
	}

	// Render the graphic into it
	decodeGraphicToImage(visualObjectMap[1])

	// Save it
	pngFile, _ := os.Create("image.png")
	png.Encode(pngFile, img)

	// printing binary!
	// dave := data[10535]
	// fmt.Println(strconv.FormatInt(int64(dave), 2))
}
