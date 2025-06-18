package main

import (
	"fmt"
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
			fmt.Fprintf(os.Stderr, "Apply command failed: %v\n", err)
			os.Exit(1)
		}
	default:
		u.PrintUsage()
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
