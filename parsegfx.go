package main

import (
	"fmt"
	"io/ioutil"
	"os"
	//"strconv"
)

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

func main() {

	if len(os.Args) != 3 {
		//fmt.Println("Usage: parsegfx <filename>")
		//return
	}

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

	// 8 pixel lookup tables
	pixTableOffs := 7620
	pixTableSize := []int{16,16,16,4,4,4,4,4}

	for i, e := range pixTableSize {
		fmt.Println("Pixel table",i+1)
		printPixelTable(data[pixTableOffs:pixTableOffs+e])
		pixTableOffs += e
		fmt.Println()
	}

	ptr := 8006 // lookup table offset

	for i := 0; i < 59; i++ {
		tableAddr := (uint16(data[ptr+1]) << 8) | uint16(data[ptr+0])
		fileOffst := tableAddr + 4096 - 6400
		if (i>0) {
			outputStr := fmt.Sprintf("Object %02d : Data %04x (File offset %04d), Width %02d bytes (%02d pixels), Height %02d, Extra %02x", i, tableAddr, fileOffst, data[ptr+2], data[ptr+2]*2, data[ptr+3], data[ptr+4])
			fmt.Println(outputStr)
		}
	    ptr += 5
	}

	// printing binary!
	// dave := data[10535]
	// fmt.Println(strconv.FormatInt(int64(dave), 2))
}
