package main

import (
	"fmt"
	"io/ioutil"
	"os"
//	"strconv"
)

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

	var totalBytes = len(data)
	fmt.Println(totalBytes, "total bytes")

	ptr := 8006 // lookup table offset
	
	// todo add absolute file addresses
	for i := 0; i < 59-50; i++ {
		tableAddr := (uint16(data[ptr+1]) << 8) | uint16(data[ptr+0])
		fileOffst := tableAddr + 4096 - 6400
		outputStr := fmt.Sprintf("Object %02d : Data %04x (File offset %04d), Width %02d bytes (%02d pixels), Height %02d, Extra %02x", i, tableAddr, fileOffst, data[ptr+2], data[ptr+2]*2, data[ptr+3], data[ptr+4])
		fmt.Println(outputStr)
	    ptr += 5
	}
}
