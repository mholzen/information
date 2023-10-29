package triples

import (
	"reflect"
	"runtime"
)

type NodeBoolFunctionNode NodeBoolFunction

// things you can do with a function:
// - evaluate its output
// - view its code
// - time it took to run
// - memory it used
// - number of instructions

type NodeBoolFunction func(Node) bool

func (n NodeBoolFunctionNode) String() string {
	return runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
}

func (n NodeBoolFunctionNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case UnaryFunctionNode:
		return n.String() < other.String()
	default:
		// TODO: variable ordering
		return false
	}
}

func NewNodeMatchAny1() NodeBoolFunctionNode {
	return func(node Node) bool {
		return true
	}
}

func NewNodeMatchAnyNumber1() NodeBoolFunctionNode {
	return func(node Node) bool {
		_, ok := node.(NumberNode)
		return ok
	}
}

func NewNodeMatchAnyIndex1() NodeBoolFunctionNode {
	return func(node Node) bool {
		_, ok := node.(IndexNode)
		return ok
	}
}

func NewNodeMatchAnyAnonymous1() NodeBoolFunctionNode {
	return func(node Node) bool {
		_, ok := node.(AnonymousNode)
		return ok
	}
}
