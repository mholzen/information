package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// func Test_TripleMatch(t *testing.T) {
// 	tpl, err := tr.NewTriple("a", tr.TypeNode, "StringNode")
// 	require.Nil(t, err)
// 	match := NewTripleFMatch(tpl)

// 	toMatch, err := tr.NewTriple("a", "type.result", "StringNode")
// 	require.Nil(t, err)
// 	x, err := match(toMatch)
// 	require.Nil(t, err)

// 	assert.True(t, x)
// }

func Test_Contains(t *testing.T) {
	tpls := tr.NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("b", "c", 2)
	tpls.AddTriple("e", "f", 3)

	toFind, err := tr.NewTriple("a", tr.NewNodeMatchAny(), 1)
	require.Nil(t, err)
	err = NewContains(toFind, tpls)(tpls)
	require.Nil(t, err)

	res, err := tpls.Map(NewReferences())
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
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("b", "c", 2)
	tpls.AddTriple("e", "f", 3)

	toFind := tr.NewTriples()
	toFind.AddTriple("a", tr.NewNodeMatchAny(), 1)
	toFind.AddTriple("b", tr.NewNodeMatchAny(), 2)

	res, err := NewContainsTriples(toFind)(tpls)
	require.Nil(t, err)

	// res, err := tpls.Map(NewReferences())
	// require.Nil(t, err)

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
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("b", "c", 2)
	tpls.AddTriple("e", "f", 3)

	toFind := tr.NewTriples()
	toFind.AddTriple("a", tr.NewNodeMatchAny(), 1)
	toFind.AddTriple("b", tr.NewNodeMatchAny(), 2)

	res, err := NewContainsTriples(toFind)(tpls)
	require.Nil(t, err)

	// res, err := tpls.Map(NewReferences())
	// require.Nil(t, err)

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
	tpls.AddTriple(tr.NewAnonymousNode(), "name", "Marc")
	tpls.AddTriple(tr.NewAnonymousNode(), "length", NewVariableNode())

	query, _ := tr.NewTriple("Marc", "length", 4)

	functions := tr.NewTriples()
	functions.AddTriple(tr.LengthFunction, ComputeNode, "length")

	res, err := NewContainsOrComputeMapper(query, functions)(tpls)

	require.Nil(t, err)
	assert.Len(t, res.TripleSet, 3)
}
