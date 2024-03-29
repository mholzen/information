package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UnaryFunctionName(t *testing.T) {
	assert.Equal(t, "TypeFunction", TypeFunctionNode.String())
}

func Test_UnaryFunctionNode_LessThan(t *testing.T) {
	assert.True(t, LengthFunctionNode.LessThan(TypeFunctionNode))
	assert.False(t, TypeFunctionNode.LessThan(LengthFunctionNode))
	assert.False(t, TypeFunctionNode.LessThan(TypeFunctionNode))

	a := NewAnonymousNode()
	compareFunctionAgainstAnon := TypeFunctionNode.LessThan(a)
	compareAnonAgainstFunction := a.LessThan(TypeFunctionNode)

	assert.NotEqual(t, compareFunctionAgainstAnon, compareAnonAgainstFunction)
}
