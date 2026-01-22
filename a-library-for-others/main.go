package main

import (
	"fmt"
	"io"
	"os"

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

	line, err := csvparser.ReadLine(file)
	if err != nil {
		fmt.Println("Error reading first line:", err)
		return
	}

	expectedFields := csvparser.GetNumberOfFields()

	fmt.Printf("Line: %s\n", line)
	for i := 0; i < expectedFields; i++ {
		field, err := csvparser.GetField(i)
		if err != nil {
			fmt.Printf("Error getting field %d: %v\n", i+1, err)
			continue
		}
		fmt.Printf("Field %d: %s\n", i+1, field)
	}

	for {
		line, err := csvparser.ReadLine(file)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading line: ", err)
			return
		}

		if line != "" {
			fmt.Printf("Line: %s\n", line)

			actualFields := csvparser.GetNumberOfFields()
			if actualFields != expectedFields {
				err := parser.ErrFieldCount
				fmt.Printf("Error %v (expected %d, found %d)\n", err, expectedFields, actualFields)
			}

			for i := 0; i < expectedFields; i++ {
				field, err := csvparser.GetField(i)
				if err != nil {
					fmt.Printf("Error getting field %d: %v\n", i, err)
					continue
				}
				fmt.Printf("Field %d: %s\n", i+1, field)
			}
		}
	}
}
