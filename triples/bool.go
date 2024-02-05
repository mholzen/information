package triples

import (
	"fmt"
	"reflect"
	"runtime"
	"strings"
)

type BoolValue bool

func (b BoolValue) String() string {
	return fmt.Sprintf("%t", b)
}

func (b BoolValue) Compare(other BoolValue) int {
	if b == other {
		return 0
	}
	if b {
		// sort true before false so that it sorts lexically first
		// (true "more important" than false)
		return -1
	}
	return 1
}

type BoolNode = CreatedComparableNode[BoolValue]

var boolNodes CreatedComparableNodes[BoolValue] = make(CreatedComparableNodes[BoolValue])

func NewBoolNode(b bool) BoolNode {
	return boolNodes.NewNode(BoolValue(b))
}

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
