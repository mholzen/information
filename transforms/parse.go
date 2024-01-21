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
	"strconv"
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

func NewNodeFromString(str string) t.Node {
	switch {
	case len(str) == 0:
		return t.NewStringNode(str)
	case str[0] == '?':
		return NewVariableNode()
	case str[0] == '_':
		return t.NewAnonymousNode()
	default:
		if num, err := strconv.Atoi(str); err == nil {
			return t.NewIndexNode(num)
		}
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return t.NewFloatNode(num)
		}
		return t.NewStringNode(str)
	}
}

func NewTripleFromString(triple string) (t.Triple, error) {
	atoms := strings.Split(triple, " ")
	if len(atoms) != 3 {
		return t.Triple{}, fmt.Errorf("invalid string '%s'", triple)
	}
	subject := NewNodeFromString(atoms[0])
	predicate := NewNodeFromString(atoms[1])
	object := NewNodeFromString(atoms[2])

	return t.NewTripleFromNodes(subject, predicate, object), nil
}

func NewTriplesFromStrings(triples ...string) (*t.Triples, error) {
	res := t.NewTriples()
	for _, triple := range triples {
		triple, err := NewTripleFromString(triple)
		if err != nil {
			return res, err
		}
		res.Add(triple)
	}
	return res, nil
}

func NewNamedTriples(triples ...string) (*t.Triples, error) {
	return NewNamedNodeMap().NewTriples(triples...)
}

type NamedNodeMap NodeMap

func NewNamedNodeMap() NamedNodeMap {
	return make(NamedNodeMap)
}

func NewNamedNodeMapFromTriples(names *t.Triples) NamedNodeMap {
	res := make(NamedNodeMap)
	for _, triple := range names.TripleSet {
		if triple.Predicate.String() == "name" {
			res[triple.Object.String()] = triple.Subject
		}
	}
	return res
}

func (m NamedNodeMap) NewTriples(triples ...string) (*t.Triples, error) {
	res := t.NewTriples()
	for _, triple := range triples {
		triple, err := m.NewTriple(triple)
		if err != nil {
			return res, err
		}
		res.Add(triple)
	}
	return res, nil
}

func (m NamedNodeMap) NewTriple(triple string) (t.Triple, error) {
	atoms := strings.Split(triple, " ")
	if len(atoms) != 3 {
		return t.Triple{}, fmt.Errorf("invalid string '%s'", triple)
	}
	subject, err := m.NewNode(atoms[0])
	if err != nil {
		return t.Triple{}, t.SubjectPosition.WrapError(err)
	}
	predicate, err := m.NewNode(atoms[1])
	if err != nil {
		return t.Triple{}, t.PredicatePosition.WrapError(err)
	}
	object, err := m.NewNode(atoms[2])
	if err != nil {
		return t.Triple{}, t.ObjectPosition.WrapError(err)
	}

	return t.NewTripleFromNodes(subject, predicate, object), nil
}

func (m NamedNodeMap) NewNode(str string) (t.Node, error) {
	switch {
	case len(str) == 0:
		return t.NewStringNode(str), nil
	case str[0] == '?':
		return m.GetOrSet(NewVariableNode(), str), nil
	case str[0] == '_':
		return m.GetOrSet(t.NewAnonymousNode(), str), nil
	case strings.HasSuffix(str, "()"):
		val, ok := m.Get(str[:len(str)-2])
		if !ok {
			return nil, fmt.Errorf("unknown function '%s' (dict contains %v)", str, m)
		}
		return val, nil
	default:
		if num, err := strconv.Atoi(str); err == nil {
			return t.NewIndexNode(num), nil
		}
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			return t.NewFloatNode(num), nil
		}
		return t.NewStringNode(str), nil
	}
}

func (m NamedNodeMap) GetOrSet(node t.Node, str string) t.Node {
	if len(str) > 0 {
		n, ok := m[str]
		if ok {
			return n
		} else {
			m[str] = node
			return node
		}
	} else {
		return node
	}
}

func (m NamedNodeMap) Get(str string) (t.Node, bool) {
	if len(str) > 0 {
		val, ok := m[str]
		return val, ok
	} else {
		return nil, false
	}
}
