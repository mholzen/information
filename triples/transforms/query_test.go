package transforms

import (
	"testing"

	"github.com/labstack/gommon/log"
	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VariableNode(t *testing.T) {
	x := NewVariableNode()
	var y Node = x
	_, ok := y.(VariableNode)
	assert.True(t, ok)
}

func Test_VariableList_Traverse(t *testing.T) {
	nodes := NodeList{NewStringNode("a"), NewStringNode("b"), NewStringNode("c")}
	vars := VariableList{NewVariableNode(), NewVariableNode()}
	r := vars.Traverse(nodes)
	assert.Len(t, r, 9)
}

func Test_NewQueryMapper(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("a", "b", 2)
	tpls.AddTriple("c", "d", 1)
	tpls.AddTriple("c", "d", 2)

	query := NewTriples()
	query.AddTriple("a", "b", 1)
	query.AddTriple("c", "d", 1)

	res, err := tpls.Map(NewQueryMapper(query))
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2)
	assert.True(t, refs.ContainsTriple("a", "b", 1))
	assert.True(t, refs.ContainsTriple("c", "d", 1))
}

func Test_NewQueryMapperWithMatches(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("a", "b", 2)
	tpls.AddTriple("c", "d", 1)
	tpls.AddTriple("c", "d", 2)
	tpls.AddTriple("c", "d", "c")

	query := NewTriples()
	query.AddTriple(NodeMatchAny, "d", NodeMatchAnyIndex)

	res, err := tpls.Map(NewQueryMapper(query))
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2)
	assert.True(t, refs.ContainsTriple("c", "d", 1))
	assert.True(t, refs.ContainsTriple("c", "d", 2))
	assert.False(t, refs.ContainsTriple("c", "d", "c"))
}

func Test_NewQueryMapperWithMatchesAndJoins(t *testing.T) {
	tpls := NewTriples()
	marc1 := NewAnonymousNode()
	tpls.AddTriple(marc1, "first", "john")
	tpls.AddTriple(marc1, "age", 52)

	marc2 := NewAnonymousNode()
	tpls.AddTriple(marc2, "first", "marc")
	tpls.AddTriple(marc2, "age", 16)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, "first", NodeMatchAnyString) // VariableNode matches anything but are used to join
	query.AddTriple(x, "age", 52)

	res, err := tpls.Map(NewQueryMapper(query))
	require.Nil(t, err)
	log.Debug(res)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2) // TODO: should join on both query triples -- iterate over anonymous nodes in the query
	assert.True(t, refs.ContainsTriple(marc1, "age", "52"))
}
