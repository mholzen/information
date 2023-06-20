package triples

import (
	"fmt"
	"log"
	"time"

	"storj.io/common/uuid"
)

type NodeValue[T any] interface {
	String() string
}

type CreatedNode[T NodeValue[T]] struct {
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

func NewCreatedNode[T NodeValue[T]](value T) CreatedNode[T] {
	return CreatedNode[T]{
		Value:   value,
		Created: time.Now(),
	}
}

type ComparableNode[T any] interface {
	String() string
	Compare(other T) int
}

type CreatedComparableNode[T ComparableNode[T]] struct {
	Value   T
	Created time.Time
}

func (n CreatedComparableNode[T]) String() string {
	return n.Value.String()
}

func (n CreatedComparableNode[T]) LessThan(other Node) bool {
	if same, ok := other.(CreatedNode[T]); ok {
		return n.Value.Compare(same.Value) < 0
	} else {
		return n.Created.Compare(same.Created) < 0
	}
}

func NewCreatedComparableNode[T ComparableNode[T]](value T) CreatedComparableNode[T] {
	return CreatedComparableNode[T]{
		Value:   value,
		Created: time.Now(),
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

type FloatType float64

func (i FloatType) String() string {
	return fmt.Sprintf("%f", i)
}

func (i FloatType) Compare(other FloatType) int {
	return int(float64(i) - float64(other))
}

type IndexNode = CreatedNode[Index]

func NewNode(value any) (Node, error) {
	switch typedValue := value.(type) {
	case Node:
		return typedValue, nil
	case string:
		return NewStringNode(typedValue), nil
	case int:
		return NewCreatedNode(Index(typedValue)), nil
	case float64:
		return NewCreatedNode(FloatType(typedValue)), nil
	case nil:
		return NewAnonymousNode(), nil
	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}
