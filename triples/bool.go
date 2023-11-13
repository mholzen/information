package triples

import (
	"reflect"
	"runtime"
)

// things you can do with a function:
// - evaluate its output
// - view its code
// - time it took to run
// - memory it used
// - number of instructions

type NodeBoolFunction func(Node) bool

func (n NodeBoolFunction) String() string {
	return runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
}

func (n NodeBoolFunction) LessThan(other Node) bool {
	switch other := other.(type) {
	case UnaryFunctionNode:
		return n.String() < other.String()
	default:
		// TODO: variable ordering
		return false
	}
}

func NodeMatchAny(node Node) bool {
	return true
}

func NodeMatchIndex(node Node) bool {
	_, ok := node.(IndexNode)
	return ok
}

func NodeMatchAnonymous(node Node) bool {
	_, ok := node.(AnonymousNode)
	return ok
}
