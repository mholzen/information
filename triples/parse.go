package triples

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
)

// parse examples.json into a Go data structure

func Read(filename string) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	// read the file
	file, err := os.Open(filename)
	if err != nil {
		return buffer, err
	}
	defer file.Close()

	// remove comments
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()

		// remove comments from line
		line = RemoveComment(line)

		buffer.WriteString(line + "\n")
	}
	return buffer, nil
}

func Parse(filename string) (*Triples, error) {
	buffer, err := Read(filename)
	if err != nil {
		return nil, err
	}
	return parseString(buffer)
}

func RemoveComment(line string) string {
	re := regexp.MustCompile(`(//)(?:([^"]|"[^"]*")*)$`)
	lines := re.Split(line, 2)
	return lines[0]
}

type lineCountingReader struct {
	r      io.Reader
	last   []byte
	offset int64
	lineno int
}

func (lcr *lineCountingReader) Read(p []byte) (n int, err error) {
	lcr.lineno += bytes.Count(lcr.last, []byte{'\n'})
	lcr.offset += int64(len(lcr.last))
	n, err = lcr.r.Read(p)
	lcr.last = make([]byte, n)
	copy(lcr.last, p[:n])
	return
}

func (lcr *lineCountingReader) Lineno(offset int64) int {
	offset -= lcr.offset
	return lcr.lineno + bytes.Count(lcr.last[:offset], []byte{'\n'})
}

func parseString(input bytes.Buffer) (*Triples, error) {
	lcr := &lineCountingReader{r: bufio.NewReader(&input)}

	decoder := json.NewDecoder(lcr)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		if syntaxErr, ok := err.(*json.SyntaxError); ok {
			// should use syntaxErr.Offset
			return nil, fmt.Errorf("syntax error on line %d: %s", lcr.Lineno(syntaxErr.Offset), syntaxErr)
		}
		return nil, err
	}

	res := NewTriples()

	transformer := NewParser(data)
	err = res.Transform(transformer)
	log.Printf("res: %+v", res.Nodes)
	return res, err

	// return oldParse(data, res)
}

func oldParse(data interface{}, res *Triples) (*Triples, error) {
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
	return nil, nil
}
