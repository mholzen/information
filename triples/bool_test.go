package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

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
