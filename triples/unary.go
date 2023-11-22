package triples

import (
	"fmt"
	"reflect"
	"regexp"
	"runtime"
	"strings"
)

type UnaryFunctionNode UnaryOperator

// things you can do with a function:
// - evaluate its output
// - view its code
// - time it took to run
// - memory it used
// - number of instructions

type UnaryOperator func(Node) (Node, error)

func (n UnaryFunctionNode) String() string {
	res := runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
	// trailing string after .
	return res[strings.LastIndex(res, ".")+1:]

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

var SquareFunctionNode UnaryFunctionNode = UnaryFunctionNode(Square)

var TypeFunctionNode UnaryFunctionNode = UnaryFunctionNode(func(node Node) (Node, error) {
	t := reflect.TypeOf(node).String()
	return NewStringNode(t), nil
})

var LengthFunction UnaryFunctionNode = UnaryFunctionNode(func(node Node) (Node, error) {
	if n, ok := node.(StringNode); ok {
		return NewIndexNode(len(n.String())), nil
	} else {
		return nil, fmt.Errorf("expected StringNode, got %T", node)
	}
})

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

func NewNodeMatchAnyNumber() UnaryFunctionNode {
	return func(node Node) (Node, error) {
		if _, ok := node.(NumberNode); ok {
			return NewNumberNode(1), nil
		}
		return NewNumberNode(0), nil
	}
}

func NewNodeMatchAnyIndex() UnaryFunctionNode {
	return func(node Node) (Node, error) {
		if _, ok := node.(IndexNode); ok {
			return NewNumberNode(1), nil
		}
		return NewNumberNode(0), nil
	}
}

func NewNodeMatchAnyAnonymous() UnaryFunctionNode {
	// TODO: refactor with NewNodeMatchAnyIndex
	return func(node Node) (Node, error) {
		if _, ok := node.(AnonymousNode); ok {
			return NewNumberNode(1), nil
		}
		return NewNumberNode(0), nil
	}
}

// This should be a UnaryTripleFunction
type TripleFunctionNode func(Triple) (Node, error)

func (n TripleFunctionNode) String() string {
	return runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
}

func (n TripleFunctionNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case TripleFunctionNode:
		return n.String() < other.String()
	default:
		// TODO: variable ordering
		return false
	}
}

func GetUnaryNodes(list NodeList) []UnaryFunctionNode {
	res := make([]UnaryFunctionNode, 0)
	for _, node := range list {
		if node, ok := node.(UnaryFunctionNode); ok {
			res = append(res, node)
		}
	}
	return res
}
