package triples

import (
	"fmt"
	"time"
)

type IndexValue int

func (i IndexValue) String() string {
	return fmt.Sprintf("%d", i)
}

func (i IndexValue) Compare(other IndexValue) int {
	return int(i) - int(other)
}

type IndexNode = CreatedComparableNode[IndexValue]

var indexNodes CreatedComparableNodes[IndexValue] = make(CreatedComparableNodes[IndexValue])

func NewIndexNode(value int) IndexNode {
	return indexNodes.NewNode(IndexValue(value))
}

type FloatNode = CreatedComparableNode[FloatValue]

var floatNodes CreatedComparableNodes[FloatValue] = make(CreatedComparableNodes[FloatValue])

func NewFloatNode(value float64) FloatNode {
	return floatNodes.NewNode(FloatValue(value))
}

type FloatValue float64

func (i FloatValue) String() string {
	return fmt.Sprintf("%f", i)
}

func (i FloatValue) Compare(other FloatValue) int {
	return int(float64(i) - float64(other))
}

type NumberNode struct {
	Value   float64
	Created time.Time
}

var numberNodes map[float64]NumberNode = make(map[float64]NumberNode)

func NewNumberNode(value float64) NumberNode {
	if node, ok := numberNodes[value]; ok {
		return node
	} else {
		numberNodes[value] = NumberNode{
			Value:   value,
			Created: time.Now(),
		}
		return numberNodes[value]
	}
}
func (n NumberNode) String() string {
	return fmt.Sprintf("%f", n.Value)
}

func (n NumberNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case NumberNode:
		return n.Value < other.Value
	case IndexNode:
		return n.Created.Compare(other.Created) < 0
	case AnonymousNode:
		return n.Created.Compare(other.Created) < 0
	default:
		return false
	}
}
