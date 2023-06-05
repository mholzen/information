package triples

import (
	"fmt"
	"time"
)

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
