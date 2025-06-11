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
		log.Fatal("Error: no command")
	}
	switch os.Args[1] {
	case "header":
		b.HandleHeaderCommand()
	case "apply":
		t.HandleRotateCommand()
		t.HandleApplyCommand()
	default:
		u.PrintUsage()
		log.Fatalf("Error: unknown command %s", os.Args[1])
	}
}
