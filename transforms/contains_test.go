package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Contains(t *testing.T) {
	tpls := tr.NewTriples()
	tpls.AddTripleFromAny("a", "b", 1)
	tpls.AddTripleFromAny("b", "c", 2)
	tpls.AddTripleFromAny("e", "f", 3)

	toFind, err := tr.NewTripleFromAny("a", tr.NewNodeMatchAny(), 1)
	require.Nil(t, err)
	err = NewContains(toFind, tpls)(tpls)
	require.Nil(t, err)

	res, err := tpls.Map(ReferencesMapper)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)

	assert.True(t, res.Contains(tr.Triple{
		Subject:   tr.NewStringNode("a"),
		Predicate: tr.NewStringNode("b"),
		Object:    tr.NewIndexNode(1), // IndexNode really?
	}), res.String())
}

func Test_Contains_triples(t *testing.T) {
	tpls := tr.NewTriples()
	tpls.AddTripleFromAny("a", "b", 1)
	tpls.AddTripleFromAny("b", "c", 2)
	tpls.AddTripleFromAny("e", "f", 3)

	toFind := tr.NewTriples()
	toFind.AddTripleFromAny("a", tr.NewNodeMatchAny(), 1)
	toFind.AddTripleFromAny("b", tr.NewNodeMatchAny(), 2)

	res, err := NewContainsTriples(toFind)(tpls)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 2)

	assert.True(t, res.Contains(tr.Triple{
		Subject:   tr.NewStringNode("a"),
		Predicate: tr.NewStringNode("b"),
		Object:    tr.NewIndexNode(1), // IndexNode really?
	}), res.String())

	assert.True(t, res.Contains(tr.Triple{
		Subject:   tr.NewStringNode("b"),
		Predicate: tr.NewStringNode("c"),
		Object:    tr.NewIndexNode(2), // IndexNode really?
	}), res.String())

}

func Test_Contains_triples_utility(t *testing.T) {
	tpls := tr.NewTriples()
	tpls.AddTripleFromAny("a", "b", 1)
	tpls.AddTripleFromAny("b", "c", 2)
	tpls.AddTripleFromAny("e", "f", 3)

	toFind := tr.NewTriples()
	toFind.AddTripleFromAny("a", tr.NewNodeMatchAny(), 1)
	toFind.AddTripleFromAny("b", tr.NewNodeMatchAny(), 2)

	res, err := NewContainsTriples(toFind)(tpls)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 2)

	assert.True(t, res.Contains(tr.Triple{
		Subject:   tr.NewStringNode("a"),
		Predicate: tr.NewStringNode("b"),
		Object:    tr.NewIndexNode(1), // IndexNode really?
	}), res.String())

	assert.True(t, res.Contains(tr.Triple{
		Subject:   tr.NewStringNode("b"),
		Predicate: tr.NewStringNode("c"),
		Object:    tr.NewIndexNode(2), // IndexNode really?
	}), res.String())
}

func Test_NewContainsOrComputeMapper(t *testing.T) {
	tpls := tr.NewTriples()
	tpls.AddTripleFromAny(tr.NewAnonymousNode(), "name", "Marc")
	tpls.AddTripleFromAny(tr.NewAnonymousNode(), "length", NewVariableNode())

	query, _ := tr.NewTripleFromAny("Marc", "length", 4)

	functions := tr.NewTriples()
	functions.AddTripleFromAny(tr.LengthFunctionNode, ComputeNode, "length")

	res, err := NewContainsOrComputeMapper(query, functions)(tpls)

	require.Nil(t, err)
	assert.Len(t, res.TripleSet, 3)
}
