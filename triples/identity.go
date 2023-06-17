package triples

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type Node interface {
	String() string
	LessThan(Node) bool
	// TODO: consider Created field
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
	case Node:
		return typedValue, nil
	case string:
		return NewStringNode(typedValue), nil
	case int:
		return NewStringNode(fmt.Sprintf("%d", typedValue)), nil // TODO: create NumberNode
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

func NodeString(node Node) string {
	return fmt.Sprintf("%T:%s", node, node)
}

// use guid as a unique identifier for each node
type AnonymousNode struct {
	UUID    uuid.UUID
	Created time.Time
}

func NewAnonymousNode() AnonymousNode {
	return AnonymousNode{
		UUID:    uuid.New(),
		Created: time.Now(),
	}
}
func (n AnonymousNode) String() string {
	return n.UUID.String()[0:8]
}
func (n AnonymousNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case AnonymousNode:
		return n.Created.Compare(other.Created) < 0
	case IndexNode:
		return n.Created.Compare(other.Created) < 0
	case StringNode:
		return n.Created.Compare(other.Created) < 0
	default:
		return false
	}
}

type IndexNode struct {
	Value   int
	Created time.Time
}

func (n IndexNode) String() string {
	return fmt.Sprint(int(n.Value))
}
func NewIndexNode(value int) IndexNode {
	return IndexNode{
		Value:   value,
		Created: time.Now(),
	}
}

func (n IndexNode) LessThan(other Node) bool {
	switch other := other.(type) {
	case IndexNode:
		return n.Value < other.Value
	case AnonymousNode:
		return n.Created.Compare(other.Created) < 0
	case StringNode:
		return n.Created.Compare(other.Created) < 0
	default:
		return false
	}
}

func (source *Triples) NewNodeFromTriple(triple Triple) AnonymousNode {
	container := NewAnonymousNode()
	source.NewTripleFromNodes(container, Subject, triple.Subject)
	source.NewTripleFromNodes(container, Predicate, triple.Predicate)
	source.NewTripleFromNodes(container, Object, triple.Object)
	return container
}

func (source *Triples) NewNodeFromTriples(triples TripleList) AnonymousNode {
	container := NewAnonymousNode()
	for i, triple := range triples {
		node := source.NewNodeFromTriple(triple)
		source.NewTripleFromNodes(container, NewIndexNode(i), node)
	}
	return container
}
