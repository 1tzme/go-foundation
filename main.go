package main

import (
	"log"
	"os"

	b "bitmap/internal/bmp"
	t "bitmap/internal/transform"
	u "bitmap/internal/utils"
)

func main() {
	if len(os.Args) < 2 {
		u.PrintUsage()
		os.Exit(1)
	}
	switch os.Args[1] {
	case "header":
		b.HandleHeaderCommand()
	case "apply":
		if err := t.HandleApplyCommand(); err != nil {
			log.Fatalf("Apply command failed: %v", err)
		}
	default:
		u.PrintUsage()
		log.Fatalf("Unknown command: %s", os.Args[1])
	}
}
