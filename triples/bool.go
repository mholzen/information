package triples

import (
	"reflect"
	"runtime"
	"strings"
)

// things you can do with a function:
// - evaluate its output
// - view its code
// - time it took to run
// - memory it used
// - number of instructions

type NodeBoolFunction func(Node) bool

func (n NodeBoolFunction) String() string {
	res := runtime.FuncForPC(reflect.ValueOf(n).Pointer()).Name()
	// trailing string after .
	return res[strings.LastIndex(res, ".")+1:]
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

func NodeMatchAnyString(node Node) bool {
	_, ok := node.(StringNode)
	return ok
}

func NodeMatchAnyIndex(node Node) bool {
	_, ok := node.(IndexNode)
	return ok
}

func NodeMatchAnyAnonymous(node Node) bool {
	_, ok := node.(AnonymousNode)
	return ok
}
