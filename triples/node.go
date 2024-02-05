package triples

import (
	"fmt"
	"reflect"
	"strings"
	"time"
)

type Node interface {
	String() string
	LessThan(Node) bool
}

func NodeEquals(v1, v2 Node) bool {
	return reflect.TypeOf(v1) == reflect.TypeOf(v2) &&
		v1.String() == v2.String()
}

func NodeLessThan(v1, v2 Node) bool {
	t := strings.Compare(reflect.TypeOf(v1).Name(), reflect.TypeOf(v2).Name())
	switch {
	case t < 0:
		return true
	case t > 0:
		return false
	default:
		return v1.LessThan(v2)
	}
}

//
// Nodes compared using their creation time
//

type Value[T any] interface {
	String() string
}

type CreatedNode[T Value[T]] struct {
	Value   T
	Created time.Time
}

func (n CreatedNode[T]) String() string {
	return n.Value.String()
}

func (n CreatedNode[T]) LessThan(other Node) bool {
	if same, ok := other.(CreatedNode[T]); ok {
		return n.Created.Compare(same.Created) < 0
	} else {
		return NodeLessThan(n, other)
	}
}

func (s CreatedNode[T]) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"created": "%s",
		"value": "%s",
		"type": "%T"
	}`, s.Created.Format(time.RFC3339), s.Value.String(), s.Value)), nil
}

func NewCreatedNode[T Value[T]](value T) CreatedNode[T] {
	return CreatedNode[T]{
		Value:   value,
		Created: time.Now(),
	}
}

//
// Nodes that have comparable values with eachother but
// compared against other nodes user their creation time
//

type ComparableValue[T Value[T]] interface {
	String() string
	Compare(other T) int
}

type CreatedComparableNode[T ComparableValue[T]] struct {
	Value   T
	Created time.Time
}

func (s CreatedComparableNode[T]) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`{
		"created": "%s",
		"value": "%s",
		"type": "%T"
	}`, s.Created.Format(time.RFC3339), s.Value.String(), s.Value)), nil
}

func (n CreatedComparableNode[T]) String() string {
	return n.Value.String()
}

func (n CreatedComparableNode[T]) LessThan(other Node) bool {
	switch typedValue := other.(type) {
	case CreatedComparableNode[T]:
		return n.Value.Compare(typedValue.Value) < 0
	case CreatedNode[Value[any]]:
		// Unfortunately, this does not match any CreatedNode types
		return n.Created.Compare(typedValue.Created) < 0
	default:
		return NodeLessThan(n, other)
	}
}

func (n CreatedComparableNode[T]) Compare(other Node) int {
	if same, ok := other.(CreatedNode[T]); ok {
		return n.Value.Compare(same.Value)
	} else {
		return n.Created.Compare(same.Created)
	}
}

type CreatedComparableNodes[T ComparableValue[T]] map[string]CreatedComparableNode[T]

func (nodes CreatedComparableNodes[T]) NewNode(value T) CreatedComparableNode[T] {
	if node, ok := nodes[value.String()]; ok {
		return node
	} else {
		node := CreatedComparableNode[T]{
			Value:   value,
			Created: time.Now(),
		}
		nodes[value.String()] = node
		return node
	}
}

func NewNode(value any) (Node, error) {
	switch typedValue := value.(type) {
	case NodeBoolFunction:
		return typedValue, nil
	case func(Node) bool:
		return NodeBoolFunction(typedValue), nil
	case Node:
		return typedValue, nil
	case string:
		return NewStringNode(typedValue), nil
	case int:
		return NewIndexNode(typedValue), nil
	case float64:
		return NewFloatNode(typedValue), nil
	case bool:
		return NewBoolNode(typedValue), nil
	case nil:
		return NewAnonymousNode(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
