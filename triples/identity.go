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

func (source *Triples) NewNode(value interface{}) (Node, error) {
	return NewNode(value)
}

func NodeString(node Node) string {
	return fmt.Sprintf("%T:%s", node, node)
}

// use guid as a unique identifier for each node
type AnonymousNode1 struct {
	UUID    uuid.UUID
	Created time.Time
}

func (n AnonymousNode1) String() string {
	return n.UUID.String()[0:8]
}
func (n AnonymousNode1) LessThan(other Node) bool {
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
