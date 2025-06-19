package main

import (
	"log"
	"os"
	"strings"

	b "bitmap/internal/bmp"
	t "bitmap/internal/transform"
	u "bitmap/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		u.PrintUsage()
		log.Fatal("Error: no command")
	}
	switch os.Args[1] {
	case "header":
		b.HandleHeaderCommand()
	case "apply":
		hasRotate := false
		hasCrop := false
		for _, arg := range os.Args[2:] {
			if strings.HasPrefix(arg, "--rotate") {
				hasRotate = true
			}
			if strings.HasPrefix(arg, "--crop") {
				hasCrop = true
			}
		}

		if hasRotate && hasCrop {
			u.PrintApplyUsage()
			log.Fatal("Error: cannot combine --rotate and --crop")
		} else if hasRotate {
			t.HandleRotateCommand()
		} else if hasCrop {
			t.HandleCropCommand()
		} else {
			u.PrintApplyUsage()
			log.Fatal("Error: must specify --rotate or --crop")
		}
	default:
		u.PrintUsage()
		log.Fatalf("Error: unknown command %s", os.Args[1])
	}
}
