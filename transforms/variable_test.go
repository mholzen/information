package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VariableNode(t *testing.T) {
	x := NewVariableNode()
	var y tr.Node = x
	_, ok := y.(VariableNode)
	assert.True(t, ok)
}

func Test_VariableList_Traverse(t *testing.T) {
	nodes := tr.NodeList{tr.NewStringNode("a"), tr.NewStringNode("b"), tr.NewStringNode("c")}
	vars := VariableList{NewVariableNode(), NewVariableNode()}
	r := vars.Traverse(nodes)
	assert.Len(t, r, 9)
}

func Test_VariableMap_Set(t *testing.T) {
	triples := tr.NewTriples()
	v1 := NewVariableNode()
	v2 := NewVariableNode()
	v3 := NewVariableNode()
	triples.AddTriple(v1, v2, v3)
	variables := NewVariableMapFromTripleList(triples.GetTripleList())

	err := variables.TestOrSet(v1, tr.NewStringNode("a"))
	require.Nil(t, err)
	err = variables.TestOrSet(v1, tr.NewStringNode("a"))
	require.Nil(t, err)
	err = variables.TestOrSet(v1, tr.NewStringNode("b"))
	require.NotNil(t, err)

	variables.Clear()
	err = variables.TestOrSet(v1, tr.NewStringNode("b"))
	require.Nil(t, err)
	err = variables.TestOrSet(v1, tr.NewStringNode("a"))
	require.NotNil(t, err)
}

func Test_VariableMap_Get(t *testing.T) {
	v1 := NewVariableNode()
	v2 := NewVariableNode()
	variables := NewVariableMap(VariableList{v1, v2})

	err := variables.TestOrSet(v1, tr.Str("a"))
	require.Nil(t, err)
	err = variables.TestOrSet(v2, tr.Str("b"))
	require.Nil(t, err)

	v, err := variables.Get(v1)
	require.Nil(t, err)
	assert.Equal(t, tr.NewStringNode("a"), v)

	v, err = variables.Get(v2)
	require.Nil(t, err)
	assert.Equal(t, tr.NewStringNode("b"), v)

	_, err = variables.Get(Var())
	require.NotNil(t, err)
}

func Test_VariableMap_GetTriple(t *testing.T) {
	v1 := NewVariableNode()
	v2 := NewVariableNode()
	variables := NewVariableMap(VariableList{v1, v2})

	err := variables.TestOrSet(v1, tr.Str("a"))
	require.Nil(t, err)
	err = variables.TestOrSet(v2, tr.Str("b"))
	require.Nil(t, err)

	triple := tr.NewTripleFromNodes(v1, v2, tr.Str("c"))
	v, err := variables.GetTriple(triple)
	require.Nil(t, err)
	assert.Equal(t, tr.NewStringNode("a"), v.Subject)
	assert.Equal(t, tr.NewStringNode("b"), v.Predicate)
	assert.Equal(t, tr.NewStringNode("c"), v.Object)

	triple.Subject = Var()
	v, err = variables.GetTriple(triple)
	require.NotNil(t, err)
}
