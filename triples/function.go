package triples

import "time"

type NodeFunction func(Node) bool

type FunctionNode struct {
	Value   NodeFunction
	Name    string
	Created time.Time
}

func NewFunctionNode(value NodeFunction, name string) FunctionNode {
	return FunctionNode{
		Value:   value,
		Name:    name,
		Created: time.Now(),
	}

}

func (n FunctionNode) String() string {
	return n.Name
}

func (n FunctionNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case FunctionNode:
		return n.String() < other.String()
	case IndexNode:
		return n.Created.Compare(other.Created) < 0
	case AnonymousNode:
		return n.Created.Compare(other.Created) < 0
	default:
		return false
	}
}
