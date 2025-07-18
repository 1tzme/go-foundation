package main

import (
	"fmt"
	"os"

	"own-redis/internal/flags"
	"own-redis/internal/server"
)

func main() {
	port := flags.FlagInit()

	server := server.NewServer(port)
	err := server.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
		os.Exit(1)
	}
}
