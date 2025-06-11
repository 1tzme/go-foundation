package bmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"

	u "bitmap/intermal/utils"
)

type Header struct {
	FileType       string
	FileSize       uint32
	HeaderSize     uint32
	DibHeaderSize  uint32
	WidthInPixels  int32
	HeightInPixels int32
	PixelSize      uint16
	ImageSize      uint32
}

func NewHeader() *Header {
	return &Header{}
}

const (
	bitmapFileHeaderSize = 14
	bitmapInfoHeaderSize = 40
)

func HandleHeaderCommand() {
	if len(os.Args) < 3 {
		u.PrintHeaderUsage()
		log.Fatal("Error: no input file")
	}
	if len(os.Args) > 3 {
		u.PrintHeaderUsage()
		log.Fatal("Error: too many arguments")
	}
	inputFile := os.Args[2]
	printHeader(inputFile)
}

func printHeader(path string) {
	header := readHeader(path)

	fmt.Println("BMP Header:")
	fmt.Printf("- FileType %s\n", header.FileType)
	fmt.Printf("- FileSizeInBytes %d\n", header.FileSize)
	fmt.Printf("- Header Size %d\n", header.HeaderSize+header.DibHeaderSize)
	fmt.Println("DIB Header:")
	fmt.Printf("- DibHeaderSize %d\n", header.DibHeaderSize)
	fmt.Printf("- WidthInPixels %d\n", header.WidthInPixels)
	fmt.Printf("- HeightInPixels %d\n", header.HeightInPixels)
	fmt.Printf("- PixelSizeInBits %d\n", header.PixelSize)
	fmt.Printf("- ImageSizeInBytes %d\n", header.ImageSize)
}

func readHeader(path string) *Header {
	file, err := os.Open(path)
	if err != nil {
		log.Fatalln("Failer to open file: ", err)
	}
	defer file.Close()

	header := NewHeader()
	header.HeaderSize = bitmapFileHeaderSize

	fileType := make([]byte, 2)
	_, err = file.Read(fileType)
	if err != nil {
		log.Fatal("Failed to read file type: ", err)
	}
	header.FileType = string(fileType)
	if header.FileType != "BM" {
		log.Fatal(path, " is not a valid BMP file")
	}

	err = binary.Read(file, binary.LittleEndian, &header.FileSize)
	if err != nil {
		log.Fatal("Failed to read file size: ", err)
	}

	_, err = file.Seek(8, io.SeekCurrent)
	if err != nil {
		log.Fatal("Failed to skip reseerved and offset fields: ", err)
	}

	err = binary.Read(file, binary.LittleEndian, &header.DibHeaderSize)
	if err != nil {
		log.Fatal("Failed to read DIB header size: ", err)
	}
	if header.DibHeaderSize != bitmapInfoHeaderSize {
		log.Fatal("Unsupported DIP header size: ", header.DibHeaderSize)
	}

	err = binary.Read(file, binary.LittleEndian, &header.WidthInPixels)
	if err != nil {
		log.Fatal("Failed to read width: ", err)
	}

	err = binary.Read(file, binary.LittleEndian, &header.HeightInPixels)
	if err != nil {
		log.Fatal("Failed to read height: ", err)
	}

	_, err = file.Seek(2, io.SeekCurrent)
	if err != nil {
		log.Fatal("Failed to skip planes: ", err)
	}

	err = binary.Read(file, binary.LittleEndian, &header.PixelSize)
	if err != nil {
		log.Fatal("Failed to read pixel size: ", err)
	}
	if header.PixelSize != 24 {
		log.Fatal("Only 24 bit BMP files are allowed")
	}

	_, err = file.Seek(4, io.SeekCurrent)
	if err != nil {
		log.Fatal("Failed to skip compression: ", err)
	}

	err = binary.Read(file, binary.LittleEndian, &header.ImageSize)
	if err != nil {
		log.Fatal("Failed to read image size: ", err)
	}

	_, err = file.Seek(16, io.SeekCurrent)
	if err != nil {
		log.Fatal("Failed to skip remaining fields: ", err)
	}

	return header
}
