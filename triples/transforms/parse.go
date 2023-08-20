package transforms

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"

	. "github.com/mholzen/information/triples"
)

type Parser struct {
	*Triples
}

func (source *Parser) ParseMap(m map[string]interface{}) (Node, error) {
	container := NewAnonymousNode()
	for key, val := range m {
		err := source.ParseAdd(container, key, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *Parser) ParseSlice(slice []interface{}) (Node, error) {
	container := NewAnonymousNode()
	for i, val := range slice {
		err := source.ParseAdd(container, i, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *Parser) ParseAdd(subject Node, predicate, object any) error {
	object, err := source.Parse(object)
	if err != nil {
		return err
	}
	_, err = source.AddTriple(subject, predicate, object)
	if err != nil {
		return err
	}
	return nil
}

func (source *Parser) Parse(data any) (Node, error) {
	switch data := data.(type) {
	case float64, string:
		return source.NewNode(data)
	case map[string]interface{}:
		return source.ParseMap(data)
	case []interface{}:
		return source.ParseSlice(data)
	case [][]string:
		array := NewAnonymousNode()
		for i, stringArray := range data {
			row := NewAnonymousNode()
			for j, val := range stringArray {
				_, err := source.AddTriple(row, j, val)
				if err != nil {
					return array, err
				}
			}
			_, err := source.AddTriple(array, i, row)
			if err != nil {
				return array, err
			}
		}
		return array, nil

	case nil:
		return source.NewNode(nil)
	default:
		return nil, fmt.Errorf("unknown type '%T'", data)
	}
}

func NewParser(data any) Transformer {
	parser := Parser{}
	return func(target *Triples) error {
		parser.Triples = target
		_, err := parser.Parse(data)
		return err
	}
}

func DecodeJson(input string) (interface{}, error) {
	buffer := bytes.NewBuffer([]byte(input))
	decoder := json.NewDecoder(strings.NewReader(buffer.String()))
	var data interface{}
	err := decoder.Decode(&data)
	return data, err
}

func NewParserFromContentType(mimeType string, data io.Reader) (*TransformerWithResult, error) {
	switch {
	case strings.HasPrefix(mimeType, "application/json"), strings.HasPrefix(mimeType, "text/plain"):
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(data)
		if err != nil {
			return nil, err
		}
		return NewJsonParser(buffer.String()), nil

	case strings.HasPrefix(mimeType, "text/csv"):
		return NewCsvParser(data), nil
	default:
		return nil, fmt.Errorf("unsupported content-type: %s", mimeType)
	}
}

func NewJsonParser(json string) *TransformerWithResult {
	transformer := TransformerWithResult{}
	transformer.Transformer = func(target *Triples) error {
		data, err := DecodeJson(json)
		if err != nil {
			log.Printf("Error decoding '%s' %s", json, err)
			return err
		}
		parser := Parser{}
		parser.Triples = target
		res, err := parser.Parse(data)
		transformer.Result = &res
		return err
	}
	return &transformer
}

func NewCsvParser(data io.Reader) *TransformerWithResult {
	transformer := TransformerWithResult{}
	transformer.Transformer = func(target *Triples) error {
		array, err := csv.NewReader(data).ReadAll()
		if err != nil {
			return err
		}
		parser := Parser{}
		parser.Triples = target
		res, err := parser.Parse(array)
		transformer.Result = &res
		return err
	}
	return &transformer
}

func NewFileJsonParser(filename string) *TransformerWithResult {
	transformer := TransformerWithResult{}
	transformer.Transformer = func(target *Triples) error {
		buffer, err := Read(filename)
		if err != nil {
			return err
		}
		data, err := DecodeJson(buffer.String())
		if err != nil {
			return err
		}

		parser := Parser{}
		parser.Triples = target
		res, err := parser.Parse(data)
		target.AddTriple(res, "filename", filename)
		transformer.Result = &res

		return err
	}
	return &transformer
}

func ReadAndStripComments(filename string) (io.Reader, error) {
	var buffer bytes.Buffer

	// read the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
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

	return bytes.NewReader(buffer.Bytes()), nil
}

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

func RemoveComment(line string) string {
	re := regexp.MustCompile(`(//)(?:([^"]|"[^"]*")*)$`)
	lines := re.Split(line, 2)
	return lines[0]
}

// type lineCountingReader struct {
// 	r      io.Reader
// 	last   []byte
// 	offset int64
// 	lineno int
// }

// func (lcr *lineCountingReader) Read(p []byte) (n int, err error) {
// 	lcr.lineno += bytes.Count(lcr.last, []byte{'\n'})
// 	lcr.offset += int64(len(lcr.last))
// 	n, err = lcr.r.Read(p)
// 	lcr.last = make([]byte, n)
// 	copy(lcr.last, p[:n])
// 	return
// }

// func (lcr *lineCountingReader) Lineno(offset int64) int {
// 	offset -= lcr.offset
// 	return lcr.lineno + bytes.Count(lcr.last[:offset], []byte{'\n'})
// }

// func parseString(input bytes.Buffer) (*Triples, error) {
// 	lcr := &lineCountingReader{r: bufio.NewReader(&input)}

// 	decoder := json.NewDecoder(lcr)
// 	var data interface{}
// 	err := decoder.Decode(&data)
// 	if err != nil {
// 		if syntaxErr, ok := err.(*json.SyntaxError); ok {
// 			// should use syntaxErr.Offset
// 			return nil, fmt.Errorf("syntax error on line %d: %s", lcr.Lineno(syntaxErr.Offset), syntaxErr)
// 		}
// 		return nil, err
// 	}

// 	res := NewTriples()

// 	transformer := NewParser(data)
// 	err = res.Transform(transformer)
// 	return res, err
// }

func NewJsonTriples(data string) (*Triples, error) {
	res := NewTriples()
	err := res.Transform(NewJsonParser(data).Transformer)
	return res, err
}
