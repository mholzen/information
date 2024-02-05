package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewNodeBool(t *testing.T) {
	trueNode := NewBoolNode(true)
	anonNode := NewAnonymousNode()
	falseNode := NewBoolNode(false)

	assert.True(t, true, trueNode)
	assert.False(t, false, falseNode)

	// should compare using values (contrary to convention, sort true before false)
	assert.True(t, trueNode.LessThan(falseNode))

	// should sort by type (anon before bool or vice versa)
	boolVersusAnon := trueNode.LessThan(anonNode)
	anonVersusBool := anonNode.LessThan(falseNode)

	assert.True(t, boolVersusAnon != anonVersusBool)
}

func Test_NodeBoolFunction_String(t *testing.T) {
	var v NodeBoolFunction = NodeMatchAny
	assert.Equal(t, "NodeMatchAny", v.String())
}

func Test_NodeBoolFunction_Comparsion(t *testing.T) {
	var v1 Node = NodeBoolFunction(NodeMatchAny)
	var v2 Node = NodeBoolFunction(NodeMatchAnyString)
	var v3 Node = NodeBoolFunction(NodeMatchAnyString)
	// assert.False(t, v1 == v2) // will panic: "comparing uncomparable type triples.NodeBoolFunction"
	assert.False(t, NodeEquals(v1, v2))
	assert.True(t, NodeEquals(v2, v3))
}
