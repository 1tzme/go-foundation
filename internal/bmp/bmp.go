package bmp

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

func ReadBMP(path string) (*BMP, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %v", err)
	}
	defer file.Close()

	header, err := readHeader(file)
	if err != nil {
		return nil, err
	}
	bmp := &BMP{
		Header: *header,
		Image:  readImage(file, *header),
	}
	return bmp, nil
}

func readImage(file *os.File, header Header) Image {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)

	_, err := file.Seek(int64(header.HeaderSize+header.DibHeaderSize), io.SeekStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to seek to image data: %v\n", err)
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "failed to read pixel data: %v\n", err)
			os.Exit(1)
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

func ReadHeaders(path string) (*BMP, *DIBHeader, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	header, err := readHeader(file)
	if err != nil {
		return nil, nil, err
	}
	bmp := &BMP{
		Header: *header,
		Image:  ReadImage(file, *header),
	}
	return bmp, nil, nil
}

func ReadImage(file *os.File, header Header) Image {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)

	_, err := file.Seek(int64(header.HeaderSize+header.DibHeaderSize), io.SeekStart)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to seek to image data: %v\n", err)
		os.Exit(1)
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
			fmt.Fprintf(os.Stderr, "failed to read pixel data: %v\n", err)
			os.Exit(1)
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
}

func writeHeader(file *os.File, header *Header) error {
	if err := binary.Write(file, binary.LittleEndian, []byte(header.FileType)); err != nil {
		return fmt.Errorf("failed to write file type: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.FileSize); err != nil {
		return fmt.Errorf("failed to write file size: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(0)); err != nil {
		return fmt.Errorf("failed to write reserved: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.HeaderSize+header.DibHeaderSize); err != nil {
		return fmt.Errorf("failed to write offset: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.DibHeaderSize); err != nil {
		return fmt.Errorf("failed to write dib header size: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.WidthInPixels); err != nil {
		return fmt.Errorf("failed to write width: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.HeightInPixels); err != nil {
		return fmt.Errorf("failed to write height: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint16(1)); err != nil {
		return fmt.Errorf("failed to write planes: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.PixelSize); err != nil {
		return fmt.Errorf("failed to write pixel size: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint32(0)); err != nil {
		return fmt.Errorf("failed to write compression: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, header.ImageSize); err != nil {
		return fmt.Errorf("failed to write image size: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(0)); err != nil {
		return fmt.Errorf("failed to write pixels per meter: %v", err)
	}
	if err := binary.Write(file, binary.LittleEndian, uint64(0)); err != nil {
		return fmt.Errorf("failed to write colors used: %v", err)
	}
	return nil
}

func writeImage(file *os.File, image *Image, header Header) error {
	width := int(header.WidthInPixels)
	height := int(header.HeightInPixels)
	rowSize := ((width*3 + 3) / 4) * 4
	pixelData := make([]byte, rowSize)

	for y := height - 1; y >= 0; y-- {
		for x := 0; x < width; x++ {
			idx := x * 3
			pixelIdx := y*width + x
			pixelData[idx] = image.Pixels[pixelIdx].B
			pixelData[idx+1] = image.Pixels[pixelIdx].G
			pixelData[idx+2] = image.Pixels[pixelIdx].R
		}
		if _, err := file.Write(pixelData); err != nil {
			return fmt.Errorf("failed to write pixel data: %v", err)
		}
	}

	writtenSize := uint32(rowSize * height)
	if header.ImageSize > writtenSize {
		padding := make([]byte, header.ImageSize-writtenSize)
		if _, err := file.Write(padding); err != nil {
			return fmt.Errorf("failed to write padding: %v", err)
		}
	}
	return nil
}
