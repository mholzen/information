package triples

import (
	"log"

	"storj.io/common/uuid"
)

type AnonymousNode CreatedNode[uuid.UUID]

func (n AnonymousNode) String() string {
	return CreatedNode[uuid.UUID](n).String()
}

func (n AnonymousNode) LessThan(other Node) bool {
	// Note: cannot use CreatedNode[uuid.UUID](n).LessThan(other)
	// because the type conversion to the generic type doesn't match
	if same, ok := other.(AnonymousNode); ok {
		return n.Created.Compare(same.Created) < 0
	} else {
		return NodeLessThan(n, other)
	}
}

func NewAnonymousNode() AnonymousNode {
	value, err := uuid.New()
	if err != nil {
		log.Fatal(err)
	}
	return AnonymousNode(NewCreatedNode(value))
}
