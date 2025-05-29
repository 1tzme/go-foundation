package parser

import (
	"errors"
)

var (
    ErrQuote      = errors.New("excess or missing \" in quoted-field")
    ErrFieldCount = errors.New("wrong number of fields")
)