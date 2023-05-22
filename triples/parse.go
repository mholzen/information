package triples

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// parse examples.json into a Go data structure

func Parse(filename string) (*Triples, error) {
	// read the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer

	// remove comments
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// remove comments from line
		line = strings.Split(line, "//")[0]

		buffer.WriteString(line + "\n")
	}
	return parseString(buffer)
}

type lineCountingReader struct {
	r      io.Reader
	lineno int
}

func (lcr *lineCountingReader) Read(p []byte) (n int, err error) {
	n, err = lcr.r.Read(p)
	lcr.lineno += bytes.Count(p[:n], []byte{'\n'})
	return
}

func (lcr *lineCountingReader) Lineno() int {
	return lcr.lineno + 1
}

func parseString(input bytes.Buffer) (*Triples, error) {
	lcr := &lineCountingReader{r: bufio.NewReader(&input)}

	decoder := json.NewDecoder(lcr)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			// should use syntaxErr.Offset
			return nil, fmt.Errorf("syntax error on line %d (%d): %s", lcr.Lineno(), syntaxErr.Offset, syntaxErr)
		}
		return nil, err
	}

	res := NewTriples()

	switch typedData := data.(type) {
	case map[string]interface{}:
		_, err := res.NewTriplesFromMap(typedData)
		if err != nil {
			return nil, err
		}

	case []interface{}:
		_, err := res.NewTriplesFromSlice(typedData)
		if err != nil {
			return nil, err
		}

	default:
		log.Printf("unknown")
	}

	return res, nil
}
