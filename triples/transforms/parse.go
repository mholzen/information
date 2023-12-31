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

	t "github.com/mholzen/information/triples"
)

type Parser struct {
	*t.Triples
}

func (source *Parser) ParseMap(m map[string]interface{}) (t.Node, error) {
	container := t.NewAnonymousNode()
	for key, val := range m {
		err := source.ParseAdd(container, key, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *Parser) ParseSlice(slice []interface{}) (t.Node, error) {
	container := t.NewAnonymousNode()
	for i, val := range slice {
		err := source.ParseAdd(container, i, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *Parser) ParseAdd(subject t.Node, predicate, object any) error {
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

func (source *Parser) Parse(data any) (t.Node, error) {
	switch data := data.(type) {
	case float64, string:
		return source.NewNode(data)
	case map[string]interface{}:
		return source.ParseMap(data)
	case []interface{}:
		return source.ParseSlice(data)
	case [][]string:
		array := t.NewAnonymousNode()
		for i, stringArray := range data {
			row := t.NewAnonymousNode()
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

func NewParser(data any) t.Transformer {
	parser := Parser{}
	return func(target *t.Triples) error {
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

func NewParserFromContentType(mimeType string, data io.Reader) (*t.TransformerWithResult, error) {
	switch {
	case strings.HasPrefix(mimeType, "application/json"), strings.HasPrefix(mimeType, "text/plain"):
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(data)
		if err != nil {
			return nil, err
		}
		return NewJsonParser(buffer.String()), nil

	case strings.HasPrefix(mimeType, "text/x-log"):
		return NewLinesParser(data), nil

	case strings.HasPrefix(mimeType, "text/csv"):
		return NewCsvParser(data), nil
	default:
		return nil, fmt.Errorf("unsupported content-type: %s", mimeType)
	}
}

func NewJsonParser(json string) *t.TransformerWithResult {
	transformer := t.TransformerWithResult{}
	transformer.Transformer = func(target *t.Triples) error {
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

func NewCsvParser(data io.Reader) *t.TransformerWithResult {
	transformer := t.TransformerWithResult{}
	transformer.Transformer = func(target *t.Triples) error {
		array, err := csv.NewReader(data).ReadAll()
		if err != nil {
			return err
		}
		parser := Parser{}
		parser.Triples = target
		res, err := parser.Parse(array)
		target.AddTriple(res, "source", "CsvParser")
		transformer.Result = &res
		return err
	}
	return &transformer
}

func NewLinesParser(data io.Reader) *t.TransformerWithResult {
	transformer := t.TransformerWithResult{}
	transformer.Transformer = func(target *t.Triples) error {
		array, err := readLines(data)
		if err != nil {
			return err
		}
		var container t.Node = t.NewAnonymousNode()
		for i, line := range array {
			target.AddTriple(container, i, line)
		}
		target.AddTriple(container, "source", "LinesParser")
		transformer.Result = &container
		return err
	}
	return &transformer
}

func readLines(reader io.Reader) ([]string, error) {
	scanner := bufio.NewScanner(reader)
	var lines []string

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func NewFileJsonParser(filename string) *t.TransformerWithResult {
	transformer := t.TransformerWithResult{}
	transformer.Transformer = func(target *t.Triples) error {
		buffer, err := read(filename)
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
	buffer, err := read(filename)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(buffer.Bytes()), nil
}

func read(filename string) (bytes.Buffer, error) {
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
		line = removeComment(line)

		buffer.WriteString(line + "\n")
	}
	return buffer, nil
}

func removeComment(line string) string {
	re := regexp.MustCompile(`(//)(?:([^"]|"[^"]*")*)$`)
	lines := re.Split(line, 2)
	return lines[0]
}

func NewJsonTriples(data string) (*t.Triples, error) {
	res := t.NewTriples()
	err := res.Transform(NewJsonParser(data).Transformer)
	return res, err
}
