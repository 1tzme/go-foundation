package parser

import (
	"io"
)

func (p *MyCSVParser) ReadLine(r io.Reader) (string, error) {
	p.buffer = p.buffer[:0]
	temp := [1]byte{}
	isCR := false
	inQuotes := false
	quotedField := false

	for {
		n, err := r.Read(temp[:])
		if n > 0 {
			ch := temp[0]
			if ch == '"' {
				if !inQuotes {
					inQuotes = true
					quotedField = true
				} else {
					inQuotes = false
				}
			}
			if ch == '\n' && isCR {
				isCR = false
				if !quotedField {
					break
				}
				p.buffer = append(p.buffer, '\n')
				continue
			}
			if ch == '\n' {
				if !quotedField {
					break
				}
				p.buffer = append(p.buffer, '\n')
				continue
			}
			if ch == '\r' {
				isCR = true
				if !quotedField {
					continue
				}
				p.buffer = append(p.buffer, '\r')
				continue
			}
			p.buffer = append(p.buffer, ch)
		}

		if err != nil {
			if err == io.EOF && len(p.buffer) > 0 {
				break
			}
			return "", err
		}
	}

	if len(p.buffer) == 0 {
		return "", nil
	}

	line := string(p.buffer)
	fields, err := parseCSVLine(line)
	if err != nil {
		p.fields = nil
		return "", err
	}

	p.fields = fields
	return line, nil
}

func parseCSVLine(line string) ([]string, error) {
	if len(line) == 0 {
		return []string{""}, nil
	}

	fields := []string{}
	field := []byte{}
	inQuotes := false

	for i := 0; i < len(line); i++ {
		switch line[i] {
		case '"':
			if !inQuotes {
				inQuotes = true
			} else {
				if i+1 < len(line) && line[i+1] == '"' {
					field = append(field, '"')
					i++
				} else {
					inQuotes = false
				}
			}
		case ',':
			if inQuotes {
				field = append(field, line[i])
			} else {
				fields = append(fields, string(trimSpace(field)))
				field = field[:0]
			}
		default:
			if !inQuotes && len(field) == 0 && line[i] == ' ' {
				continue
			}
			field = append(field, line[i])
		}
	}

	if inQuotes {
		return nil, ErrQuote
	}

	fields = append(fields, string(trimSpace(field)))
	return fields, nil
}

func trimSpace(s []byte) []byte {
	start := 0
	end := len(s)
	for start < end && s[start] == ' ' {
		start++
	}
	for end > start && s[end-1] == ' ' {
		end--
	}
	return s[start:end]
}
