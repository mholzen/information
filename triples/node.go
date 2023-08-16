package triples

import (
	"fmt"
	"log"
	"time"

	"storj.io/common/uuid"
)

type Node interface {
	String() string
	LessThan(Node) bool
}

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
		return n.Value.String() < other.String()
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

type CreatedNodes[T Value[T]] map[string]CreatedNode[T]

func (nodes CreatedNodes[T]) NewNode(value T) CreatedNode[T] {
	if node, ok := nodes[value.String()]; ok {
		return node
	} else {
		node := CreatedNode[T]{
			Value:   value,
			Created: time.Now(),
		}
		nodes[value.String()] = node
		return node
	}
}

//
// Comparable
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
	if same, ok := other.(CreatedComparableNode[T]); ok {
		return n.Value.Compare(same.Value) < 0
	} else {
		return n.Created.Compare(same.Created) < 0
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

type AnonymousNode = CreatedNode[uuid.UUID]

// anonumous nodes need to be compared using their creation time

func NewAnonymousNode() AnonymousNode {
	value, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	return NewCreatedNode(value)
}

type Index int

func (i Index) String() string {
	return fmt.Sprintf("%d", i)
}

func (i Index) Compare(other Index) int {
	return int(i) - int(other)
}

var indexNodes CreatedComparableNodes[Index] = make(CreatedComparableNodes[Index])

type IndexNode = CreatedComparableNode[Index]

func NewIndexNode(value int) IndexNode {
	return indexNodes.NewNode(Index(value))
}

var floatNodes CreatedComparableNodes[FloatType] = make(CreatedComparableNodes[FloatType])

type FloatNode = CreatedComparableNode[FloatType]

func NewFloatNode(value float64) FloatNode {
	return floatNodes.NewNode(FloatType(value))
}

type FloatType float64

func (i FloatType) String() string {
	return fmt.Sprintf("%f", i)
}

func (i FloatType) Compare(other FloatType) int {
	return int(float64(i) - float64(other))
}

func NewNode(value any) (Node, error) {
	switch typedValue := value.(type) {
	case Node:
		return typedValue, nil
	case string:
		return NewStringNode(typedValue), nil
	case int:
		return NewIndexNode(typedValue), nil
	case float64:
		return NewFloatNode(typedValue), nil
	case nil:
		return NewAnonymousNode(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
