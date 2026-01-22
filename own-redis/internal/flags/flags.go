package flags

import (
	"flag"
	"fmt"
	"os"

	"own-redis/internal/utils"
)

const DefaultPort = 8080

func FlagInit() int {
	help := flag.Bool("help", false, "Show help message")
	port := flag.Int("port", DefaultPort, "Port number")

	flag.Usage = utils.Usage
	flag.Parse()

	if *help {
		utils.Usage()
		os.Exit(0)
	}

	if *port < 1 || *port > 65535 {
		fmt.Fprintf(os.Stderr, "Error: Port number %d is out of valid range\n", *port)
		os.Exit(1)
	}

	return *port
}
