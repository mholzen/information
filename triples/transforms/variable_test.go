package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_VariableMap_Set(t *testing.T) {
	triples := tr.NewTriples()
	v1 := NewVariableNode()
	v2 := NewVariableNode()
	v3 := NewVariableNode()
	triples.AddTriple(v1, v2, v3)
	variables := NewVariableMap(triples.GetTripleList())

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
