package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_References(t *testing.T) {
	tpls := NewTriples()
	container := NewAnonymousNode()
	tpl, _ := NewTriple(container, "b", 1)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(container, "c", 2)
	tpls.AddTripleReference(tpl)

	res, err := tpls.Map(ReferencesMapper)
	require.Nil(t, err)

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

func Test_ReferenceTriples(t *testing.T) {
	tpls := NewTriples()
	container := NewAnonymousNode()
	tpl, _ := NewTriple(container, "b", 1)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(container, "c", 2)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple("a", "b", 3)
	tpls.Add(tpl)

	res := ReferenceTriples(tpls)

	assert.Len(t, res.TripleSet, 6)
}

func Test_RemoveReferences(t *testing.T) {
	tpls := NewTriples()
	container := NewAnonymousNode()
	tpl, _ := NewTriple(container, "b", 1)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(container, "c", 2)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple("a", "b", 3)
	tpls.Add(tpl)

	assert.Len(t, tpls.TripleSet, 7)

	res, err := tpls.Map(RemoveReferencesMapper)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)
}

func Test_ReferenceTriplesSuper(t *testing.T) {
	tpls := NewTriples()
	container := NewAnonymousNode()
	tpl, _ := NewTriple(container, "b", 1)
	ref1 := tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(container, "c", 2)
	tpls.AddTripleReference(tpl)
	tpl, _ = NewTriple(ref1, "foo", "bar")
	tpls.Add(tpl)
	tpl, _ = NewTriple("a", "b", 3)
	tpls.Add(tpl)
	assert.Len(t, tpls.TripleSet, 8)

	res := ReferenceTriplesConnected(tpls)

	assert.Len(t, res.TripleSet, 7)
}
