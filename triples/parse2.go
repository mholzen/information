package triples

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"strings"
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

func NewJsonParserOld(json string, top *Node) (Transformer, error) {
	data, err := DecodeJson(json)
	if err != nil {
		return nil, err
	}
	parser := Parser{}
	return func(target *Triples) error {
		parser.Triples = target
		res, err := parser.Parse(data)
		*top = res

		if top != nil {
			*top = res
		}
		return err
	}, nil
}

func NewJsonParser(json string) *TransformerWithResult {
	transformer := TransformerWithResult{}
	transformer.Transformer = func(target *Triples) error {
		data, err := DecodeJson(json)
		if err != nil {
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

func NewCsvParser(data string) *TransformerWithResult {
	transformer := TransformerWithResult{}
	transformer.Transformer = func(target *Triples) error {
		array, err := csv.NewReader(strings.NewReader(data)).ReadAll()
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

func NewFileJsonParser(filename string, top *Node) (Transformer, error) {
	buffer, err := Read(filename)
	if err != nil {
		return nil, err
	}
	return NewJsonParserOld(buffer.String(), top)
}
