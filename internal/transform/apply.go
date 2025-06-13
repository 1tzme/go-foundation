package transform

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	b "bitmap/internal/bmp"
	u "bitmap/internal/utils"
)

func HandleApplyCommand() error {
	fs := flag.NewFlagSet("apply", flag.ContinueOnError)
	var filters flagSlice
	var rotates []string

	fs.Var(&filters, "filter", "Filter to apply (blue, red, green, grayscale, negative, pixelate, blur)")
	fs.Func("rotate", "Rotate direction: right, 90, 180, 270, left, -90, -180, -270", func(s string) error {
		s = strings.ToLower(s)
		if !isValidRotation(s) {
			return fmt.Errorf("invalid rotate option: %s", s)
		}
		rotates = append(rotates, s)
		return nil
	})

	args := normalizeFlags(os.Args[2:])
	if err := fs.Parse(args); err != nil {
		u.PrintApplyUsage()
		return fmt.Errorf("error parsing flags: %v", err)
	}

	if len(fs.Args()) != 2 {
		u.PrintApplyUsage()
		return fmt.Errorf("expected source and output files")
	}

	if len(filters) == 0 && len(rotates) == 0 {
		u.PrintApplyUsage()
		return fmt.Errorf("no filters or rotations specified")
	}

	for _, f := range filters {
		if !isValidFilter(f) {
			return fmt.Errorf("invalid filter: %s", f)
		}
	}

	sourceFile, outputFile := fs.Args()[0], fs.Args()[1]

	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		return fmt.Errorf("source file %s does not exist", sourceFile)
	}

	if dir := filepath.Dir(outputFile); dir != "." && dir != "" {
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			return fmt.Errorf("output directory %s does not exist", dir)
		}
	}

	bmp, err := b.ReadBMP(sourceFile)
	if err != nil {
		return fmt.Errorf("failed to read BMP: %v", err)
	}

	if len(rotates) > 0 {
		bmp.Image = *applyRotations(&bmp.Image, rotates)
		bmp.Header.WidthInPixels = int32(bmp.Image.Width)
		bmp.Header.HeightInPixels = int32(bmp.Image.Height)
	}

	if len(filters) > 0 {
		applyFilters(&bmp.Image, filters)
	}

	if err := b.WriteBMP(outputFile, bmp); err != nil {
		return fmt.Errorf("failed to write BMP: %v", err)
	}
	return nil
}

func normalizeFlags(args []string) []string {
	var normalized []string
	for _, arg := range args {
		if strings.HasPrefix(arg, "---") {
			normalized = append(normalized, strings.Replace(arg, "---", "--", 1))
		} else {
			normalized = append(normalized, arg)
		}
	}
	return normalized
}

type flagSlice []string

func (f *flagSlice) String() string {
	return ""
}

func (f *flagSlice) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func isValidFilter(filter string) bool {
	validFilters := map[string]bool{
		"blue":      true,
		"red":       true,
		"green":     true,
		"grayscale": true,
		"negative":  true,
		"pixelate":  true,
		"blur":      true,
	}
	return validFilters[filter]
}
