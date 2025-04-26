package csvparser

import (
	"errors"
	"io"
)

type CSVParser interface  {
    ReadLine(r io.Reader) (string, error)
    GetField(n int) (string, error)
    GetNumberOfFields() int
}

type CSVParserData struct {
	line string
	fields []string
	eatenSlashR bool
	prevChar byte
}

var (
    ErrQuote      = errors.New("excess or missing \" in quoted-field")
    ErrFieldCount = errors.New("wrong number of fields")
)

func (c *CSVParserData) ReadLine(r io.Reader) (string, error) {
	var buffer []byte
	var insideQuotes bool
	for {
		if c.eatenSlashR {
			buffer = append(buffer, c.prevChar)
			c.eatenSlashR = false
		}
		temp := make([]byte, 1)
		_, err := r.Read(temp)
		if err != nil {
			if err == io.EOF {
				if insideQuotes {
					return "", ErrQuote
				} else if len(buffer) > 0 {
					break
				}
				return "", io.EOF
			}
			return "", err
		}
		char := temp[0]
		if char == '"' {
			insideQuotes = !insideQuotes
		}

		if !insideQuotes {
			if char == '\r' {
				temp = make([]byte, 1)
				_, err := r.Read(temp)
				if err != nil {
					if err == io.EOF {
						if len(buffer) > 0 {
							break
						}
						return "", io.EOF
					}
					return "", err
				}
				if temp[0] == '\n' {
					break
				} else {
					c.prevChar = temp[0]
					c.eatenSlashR = true
					break
				}
			} else if char == '\n' {
				break
			}
		}

		buffer = append(buffer, char)
		c.prevChar = char
	} 
	line := fixQuote(string(buffer))
	c.fields = separateLine(line)
	c.line = line
	
	return line, nil
}

func fixQuote(str string) string {
	newStr := ""
	var i int
	for i = 0; i < len(str)-1; i++ {
		if str[i] == '"' && str[i+1] == '"' {
			continue
		} 
		newStr += string(str[i])
	}
	newStr += string(str[i])
	return newStr
}

func separateLine(line string) []string {
	fields := []string{}
	var openQuotes bool
	tempStr := ""
	for i := 0; i < len(line); i++{
		switch line[i] {
		case '"':
			openQuotes = !openQuotes
			tempStr += string(line[i])
		case ',':
			if openQuotes {
				tempStr += string(line[i])
			} else {
				fields = append(fields, tempStr)
				tempStr = ""
			}
		default:
			tempStr += string(line[i])
		}
	}
	if tempStr != "" {
		fields = append(fields, tempStr)
	}
	return fields
}