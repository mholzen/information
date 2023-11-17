package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_Variable_Set(t *testing.T) {
	x := NewVariableNode()
	require.Nil(t, x.TestOrSet(tr.NewStringNode("a")))
	require.Nil(t, x.TestOrSet(tr.NewStringNode("a")))
	require.NotNil(t, x.TestOrSet(tr.NewStringNode("b")))
}

func Test_VariableList_Clear(t *testing.T) {
	vars := VariableList{NewVariableNode(), NewVariableNode()}
	vars[0].TestOrSet(tr.NewStringNode("a"))
	vars[1].TestOrSet(tr.NewStringNode("b"))
	vars.Clear()
	require.Nil(t, vars[0].Value)
	require.Nil(t, vars[1].Value)
}
