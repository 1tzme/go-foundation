package utils

import "fmt"

func Usage() {
	fmt.Println(`Own Redis

Usage:
  own-redis [--port <N>]
  own-redis --help

Options:
  --help       Show this screen.
  --port N     Port number.`)
}
