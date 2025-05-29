package main

import (
	"fmt"
	"os"
	"io"
    "a-library-for-others/parser"
)

func main() {
    file, err := os.Open("example.csv")
    if err != nil {
        fmt.Println("Error opening file:", err)
        return
    }
    defer file.Close()

    var csvparser parser.CSVParser = MyCSVParser{}

    for {
        line, err := csvparser.ReadLine(file)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println("Error reading line:", err)
            return
        }
    }
}