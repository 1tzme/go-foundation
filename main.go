package main

import (
	"log"
	"os"

	b "bitmap/intermal/bmp"
	t "bitmap/intermal/transform"
	u "bitmap/intermal/utils"
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
	default:
		u.PrintUsage()
		log.Fatalf("Error: unknown command %s", os.Args[1])
	}
}
