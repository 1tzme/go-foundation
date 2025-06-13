package bmp

type BitmapFileHeader struct {
	FileType        [2]byte
	FileSize        uint32
	Reserved1       uint16
	Reserved2       uint16
	PixelDataOffset uint32
}

type DIBHeader struct {
	HeaderSize      uint32
	Width           int32
	Height          int32
	Planes          uint16
	BitsPerPixel    uint16
	Compression     uint32
	ImageSize       uint32
	XPixelsPerMeter int32
	YPixelsPerMeter int32
	ColorsUsed      uint32
	ImportantColors uint32
}

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
