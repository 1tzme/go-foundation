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

    var csvparser parser.CSVParser = &parser.MyCSVParser{}

    for {
        line, err := csvparser.ReadLine(file)
        if err != nil {
            if err == io.EOF {
                break
            }
            fmt.Println("Error reading line: ", err)
            return
        }

        fmt.Printf("Line: %s\n", line)
        for i := 0; i < csvparser.GetNumberOfFields(); i++ {
            field, err := csvparser.GetField(i)
            if err != nil {
                fmt.Printf("Error getting field %d: %v\n", i, err)
                continue
            }
            fmt.Printf("Field %d: %s\n", i, field)
        }
    }
}