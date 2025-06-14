package transform

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	b "bitmap/internal/bmp"
	u "bitmap/internal/utils"
)

func HandleRotateCommand() {
	var rotates []string
	rotateFlag := flag.NewFlagSet("apply", flag.ExitOnError)
	rotateFlag.Func("rotate", "Rotate directIion: right, 90, 180, 270, left, -90, -180, -270", func(s string) error {
		s = strings.ToLower(s)
		if !isValidRotation(s) {
			return fmt.Errorf("invalid rotate option: %s", s)
		}
		rotates = append(rotates, s)
		return nil
	})

	if err := rotateFlag.Parse(os.Args[2:]); err != nil {
		log.Fatalf("Error parsing flags: %v", err)
	}

	if len(rotateFlag.Args()) != 2 {
		u.PrintApplyUsage()
		log.Fatal("Error: expected source and output file")
	}

	sourceFile, outputFile := rotateFlag.Args()[0], rotateFlag.Args()[1]
	img := b.ReadImage(sourceFile)
	rotatedImg := applyRotations(img, rotates)
	b.SaveImage(rotatedImg, outputFile)
}

func isValidRotation(option string) bool {
	validOptions := map[string]bool{
		"right": true, "90": true, "180": true, "270": true,
		"left": true, "-90": true, "-180": true, "-270": true,
	}
	return validOptions[strings.ToLower(option)]
}

func applyRotations(img *b.Image, rotations []string) *b.Image {
	result := img
	for _, rot := range rotations {
		angle := normalizeRotation(rot)
		result = rotateImage(result, angle)
	}
	return result
}

func normalizeRotation(rot string) int {
	rot = strings.ToLower(rot)
	switch rot {
	case "right", "90":
		return 90
	case "left", "-90":
		return -90
	case "180", "-180":
		return 180
	case "270", "-270":
		return 270
	default:
		log.Fatalf("Invalid rotation: %s", rot)
		return 0
	}
}

func rotateImage(img *b.Image, angle int) *b.Image {
	angle = (angle%360 + 360) % 360 // Нормализация угла
	width, height := img.Width, img.Height
	var newWidth, newHeight int
	var newPixels []b.Pixel

	switch angle {
	case 90, 270:
		newWidth, newHeight = height, width
	default:
		newWidth, newHeight = width, height
	}

	newPixels = make([]b.Pixel, newWidth*newHeight)

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			idx := y*width + x
			var newX, newY int
			switch angle {
			case 90:
				newX, newY = y, newHeight-1-x
			case 180:
				newX, newY = newWidth-1-x, newHeight-1-y
			case 270:
				newX, newY = newWidth-1-y, x
			default: // 0 градусов
				newX, newY = x, y
			}
			newIdx := newY*newWidth + newX
			newPixels[newIdx] = img.Pixels[idx]
		}
	}

	return &b.Image{
		Width:  newWidth,
		Height: newHeight,
		Pixels: newPixels,
	}
}
