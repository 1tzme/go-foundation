package main

import (
	"bitmap/intermal/bmp"
	"fmt"
	"log"
)

func main() {
	fileHeader, dibHeader, err := bmp.ReadHeaders("sample.bmp")
	if err != nil {
		log.Fatalf("Error: %v\n", err)
	}

	fmt.Println("BMP Header:")
	fmt.Printf("- FileType %s\n", fileHeader.FileType)
	fmt.Printf("- FileSizeInBytes %d\n", fileHeader.FileSize)
	fmt.Printf("- HeaderSize %d\n", fileHeader.PixelDataOffset)

	fmt.Println("DIB Header:")
	fmt.Printf("- DibHeaderSize %d\n", dibHeader.HeaderSize)
	fmt.Printf("- WidthInPixels %d\n", dibHeader.Width)
	fmt.Printf("- HeightInPixels %d\n", dibHeader.Height)
	fmt.Printf("- PixelSizeInBits %d\n", dibHeader.BitsPerPixel)
	fmt.Printf("- ImageSizeInBytes %d\n", dibHeader.ImageSize)
}
