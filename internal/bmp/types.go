package bmp

type BMP struct {
	Header Header
	Image  Image
}

type Image struct {
	Width  int
	Height int
	Pixels []Pixel
}

type Pixel struct {
	B, G, R uint8
}

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

const (
	bitmapFileHeaderSize = 14
	bitmapInfoHeaderSize = 40
)
