package bmp

import (
	"encoding/binary"
	"fmt"
	"os"
)

func ReadHeaders(path string) (*BitmapFileHeader, *DIBHeader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	var fileHeader BitmapFileHeader
	if err := binary.Read(file, binary.LittleEndian, &fileHeader); err != nil {
		return nil, nil, fmt.Errorf("failed to read file header: %w", err)
	}

	if fileHeader.FileType != [2]byte{'B', 'M'} {
		return nil, nil, fmt.Errorf("not a BMP file: invalid signature %v", fileHeader.FileType)
	}

	var dibHeader DIBHeader
	if err := binary.Read(file, binary.LittleEndian, &dibHeader); err != nil {
		return nil, nil, fmt.Errorf("failed to read DIB header: %w", err)
	}

	return &fileHeader, &dibHeader, nil
}

// package bmp

// type BMP struct {
// 	Header Header
// 	Image Image
// }

// type Image struct {
// 	Width int
// 	Height int
// 	Pixels []Pixel
// }

// type Pixel struct {
// 	B, G, R uint8
// }
