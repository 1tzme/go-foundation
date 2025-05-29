package parser

import (
	"io"
)

type CSVParser interface  {
    ReadLine(r io.Reader) (string, error)
    GetField(n int) (string, error)
    GetNumberOfFields() int
}
