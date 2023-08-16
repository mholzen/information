package triples

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
)

type UnaryFunctionNode UnaryOperator

// things you can do with a function:
// - evaluate it
// - result
// - time it took to run
// - memory it used
// - number of instructions

type UnaryOperator func(Node) (Node, error)

func (n UnaryFunctionNode) String() string {
	return runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
}

func (n UnaryFunctionNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case UnaryFunctionNode:
		return n.String() < other.String()
	default:
		// TODO: variable ordering
		return false
	}
}

func Square(node Node) (Node, error) {
	n, ok := node.(NumberNode)
	if !ok {
		return nil, fmt.Errorf("expected NumberNode, got %T", node)
	}
	return NewNumberNode(n.Value * n.Value), nil
}

var SquareNode UnaryFunctionNode = UnaryFunctionNode(Square)

func NewStringNodeMatch(re string) UnaryFunctionNode {
	return func(node Node) (Node, error) {
		if node, ok := node.(StringNode); ok {
			if match, _ := regexp.MatchString(re, node.String()); match {
				return NewNumberNode(1), nil
			}
		}
		return NewNumberNode(0), nil
	}
}

func NewNodeMatchAny() UnaryFunctionNode {
	return func(node Node) (Node, error) {
		return NewNumberNode(1), nil
	}
}
