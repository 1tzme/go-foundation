package bmp

import (
	"encoding/binary"
	"fmt"
	"log"
	"os"

	u "bitmap/internal/utils"
)

func HandleHeaderCommand() {
	if len(os.Args) != 3 {
		u.PrintHeaderUsage()
		log.Fatal("Invalid number of arguments")
	}
	printHeader(os.Args[2])
}

func printHeader(path string) {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)
	}
	defer file.Close()

	header, err := readHeader(file)
	if err != nil {
		log.Fatalf("Failed to read header: %v", err)
	}
	fmt.Printf("BMP Header:\n")
	fmt.Printf("- FileType: %s\n", header.FileType)
	fmt.Printf("- FileSize: %d bytes\n", header.FileSize)
	fmt.Printf("- HeaderSize: %d bytes\n", header.HeaderSize+header.DibHeaderSize)
	fmt.Printf("DIB Header:\n")
	fmt.Printf("- DibHeaderSize: %d bytes\n", header.DibHeaderSize)
	fmt.Printf("- Width: %d pixels\n", header.WidthInPixels)
	fmt.Printf("- Height: %d pixels\n", header.HeightInPixels)
	fmt.Printf("- PixelSize: %d bits\n", header.PixelSize)
	fmt.Printf("- ImageSize: %d bytes\n", header.ImageSize)
}

func readHeader(file *os.File) (*Header, error) {
	header := &Header{HeaderSize: bitmapFileHeaderSize}

	fileType := make([]byte, 2)
	if _, err := file.Read(fileType); err != nil {
		return nil, fmt.Errorf("failed to read file type: %v", err)
	}
	header.FileType = string(fileType)
	if header.FileType != "BM" {
		return nil, fmt.Errorf("not a valid BMP file")
	}

	if err := binary.Read(file, binary.LittleEndian, &header.FileSize); err != nil {
		return nil, fmt.Errorf("failed to read file size: %v", err)
	}
	if _, err := file.Seek(8, 1); err != nil {
		return nil, fmt.Errorf("failed to seek reserved and offset: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &header.DibHeaderSize); err != nil {
		return nil, fmt.Errorf("failed to read dib header size: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &header.WidthInPixels); err != nil {
		return nil, fmt.Errorf("failed to read width: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &header.HeightInPixels); err != nil {
		return nil, fmt.Errorf("failed to read height: %v", err)
	}
	if _, err := file.Seek(2, 1); err != nil {
		return nil, fmt.Errorf("failed to seek planes: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &header.PixelSize); err != nil {
		return nil, fmt.Errorf("failed to read pixel size: %v", err)
	}

	if header.PixelSize != 24 {
		return nil, fmt.Errorf("only 24-bit BMP files are supported")
	}

	if _, err := file.Seek(4, 1); err != nil {
		return nil, fmt.Errorf("failed to seek compression: %v", err)
	}
	if err := binary.Read(file, binary.LittleEndian, &header.ImageSize); err != nil {
		return nil, fmt.Errorf("failed to read image size: %v", err)
	}

	return header, nil
}
