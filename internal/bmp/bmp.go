package bmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ReadHeaders(path string) (*BitmapFileHeader, *DIBHeader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	header, err := readHeader(file)
	if err != nil {
		return nil, err
	}
	bmp := &BMP{
		Header: *header,
		Image:  ReadImage(file, *header),
	}
	return bmp, nil
}

func ReadImage(file *os.File, header Header) Image {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)

	_, err := file.Seek(int64(header.HeaderSize+header.DibHeaderSize), io.SeekStart)
	if err != nil {
		panic(fmt.Errorf("failed to seek to image data: %v", err))
	}

	rowSize := ((int(header.PixelSize)*width + 31) / 32) * 4
	pixelData := make([]byte, rowSize)
	image := Image{
		Width:  width,
		Height: height,
		Pixels: make([]Pixel, width*height),
	}

	for y := height - 1; y >= 0; y-- {
		_, err := file.Read(pixelData)
		if err != nil {
			panic(fmt.Errorf("failed to read pixel data: %v", err))
		}
		for x := 0; x < width; x++ {
			idx := x * 3
			pixelIdx := y*width + x
			image.Pixels[pixelIdx] = Pixel{
				B: pixelData[idx],
				G: pixelData[idx+1],
				R: pixelData[idx+2],
			}
		}
	}
	return image
}

func WriteBMP(path string, bmp *BMP) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	updateHeader(&bmp.Header)
	if err := writeHeader(file, &bmp.Header); err != nil {
		return err
	}
	if err := writeImage(file, &bmp.Image, bmp.Header); err != nil {
		return err
	}
	return nil
}

func updateHeader(header *Header) {
	header.FileType = "BM"
	header.HeaderSize = bitmapFileHeaderSize
	header.DibHeaderSize = bitmapInfoHeaderSize
	header.PixelSize = 24

	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)
	rowSize := ((width*3 + 3) / 4) * 4
	calculatedImageSize := uint32(rowSize * height)

	if header.ImageSize < calculatedImageSize {
		header.ImageSize = calculatedImageSize
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
