package triples

import (
	"fmt"

	"github.com/google/uuid"
)

type Node interface {
	String() string
}

type NodeSet map[string]struct{}

func NewNodeSet() NodeSet {
	return make(NodeSet)
}

func (set NodeSet) Add(n Node) Node {
	if _, ok := set[n.String()]; ok {
		return n
	} else {
		set[n.String()] = struct{}{}
		return n
	}
}

func (set NodeSet) Contains(node Node) bool {
	_, ok := set[node.String()]
	return ok
}

func (set NodeSet) ContainsOrAdd(node Node) bool {
	if set.Contains(node) {
		return true
	} else {
		set.Add(node)
		return false
	}
}

func (source *Triples) NewNode(value interface{}) (Node, error) {
	switch typedValue := value.(type) {
	case string:
		return NewStringNode(typedValue), nil
	case nil:
		return NewAnonymousNode(), nil
	case []interface{}:
		triples, err := source.NewTriplesFromSlice(typedValue)
		if err != nil {
			return nil, err
		}
		return source.NewNodeFromTriples(triples), nil

	case map[string]interface{}:
		triples, err := source.NewTriplesFromMap(typedValue)
		if err != nil {
			return nil, err
		}
		return source.NewNodeFromTriples(triples), nil

	default:
		return nil, fmt.Errorf("unsupported type: %T", value)
	}
}

type StringNode string

func NewStringNode(value string) StringNode {
	if node, ok := stringNodes[value]; ok {
		return node
	} else {
		stringNodes[value] = StringNode(value)
		return stringNodes[value]
	}
}
func (n StringNode) String() string {
	return string(n)
}

var Subject = NewStringNode("1-subject")
var Predicate = NewStringNode("2-predicate")
var Object = NewStringNode("3-object")
var Contains = NewStringNode("contains")

type StringNodes map[string]StringNode

var stringNodes StringNodes = make(StringNodes)

// use guid as a unique identifier for each node
type AnonymousNode uuid.UUID

func NewAnonymousNode() AnonymousNode {
	return AnonymousNode(uuid.New())
}
func (n AnonymousNode) String() string {
	return uuid.UUID(n).String()[0:8]
}

type IndexNode int

func (n IndexNode) String() string {
	return fmt.Sprint(int(n))
}

func (source *Triples) NewNodeFromTriple(triple Triple) AnonymousNode {
	container := NewAnonymousNode()
	source.NewTriple(container, Subject, triple.Subject)
	source.NewTriple(container, Predicate, triple.Predicate)
	source.NewTriple(container, Object, triple.Object)
	return container
}

func (source *Triples) NewNodeFromTriples(triples TripleList) AnonymousNode {
	container := NewAnonymousNode()
	for i, triple := range triples {
		node := source.NewNodeFromTriple(triple)
		source.NewTriple(container, IndexNode(i), node)
	}
	return container
}
