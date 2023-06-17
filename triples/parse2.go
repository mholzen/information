package triples

import (
	"fmt"
)

type TriplesModifier struct {
	*Triples
}

func (source *TriplesModifier) ParseMap(m map[string]interface{}) (Node, error) {
	container := NewAnonymousNode()
	for key, val := range m {
		err := source.ParseAdd(container, key, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *TriplesModifier) ParseSlice(slice []interface{}) (Node, error) {
	container := NewAnonymousNode()
	for i, val := range slice {
		err := source.ParseAdd(container, i, val)
		if err != nil {
			return nil, err
		}
	}
	return container, nil
}

func (source *TriplesModifier) ParseAdd(subject Node, predicate, object any) error {
	object, err := source.Parse(object)
	if err != nil {
		return err
	}
	_, err = source.NewTriple(subject, predicate, object)
	if err != nil {
		return err
	}
	return nil
}

func (source *TriplesModifier) Parse(data any) (Node, error) {
	switch data := data.(type) {
	case float64, string:
		return source.NewNode(data)
	case map[string]interface{}:
		return source.ParseMap(data)
	case []interface{}:
		return source.ParseSlice(data)
	default:
		return nil, fmt.Errorf("unknown type '%T'", data)
	}
}

func NewParser(data any) Transformer {
	tm := TriplesModifier{}
	return func(target *Triples) error {
		tm.Triples = target
		_, err := tm.Parse(data)
		return err
	}
}
