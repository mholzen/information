package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
)

func Test_References(t *testing.T) {
	tpls := NewTriples()
	container := NewAnonymousNode()
	tpl, _ := NewTriple(container, "b", 1)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(container, "c", 2)
	tpls.AddTripleReference(tpl)

	res, err := tpls.Map(NewReferences())
	assert.Nil(t, err)

	assert.Len(t, res.TripleSet, 2)

	assert.True(t, res.Contains(Triple{
		Subject:   container,
		Predicate: NewStringNode("b"),
		Object:    NewIndexNode(1),
	}))

	assert.True(t, res.Contains(Triple{
		Subject:   container,
		Predicate: NewStringNode("c"),
		Object:    NewIndexNode(2),
	}))
}