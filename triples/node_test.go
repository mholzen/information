package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CompareNode(t *testing.T) {
	a1, _ := NewNode("a")
	a2, _ := NewNode("a")
	assert.True(t, a1 == a2)
	assert.True(t, NodeEquals(a1, a2))

	n1, _ := NewNode(1)
	n2, _ := NewNode(1)
	assert.True(t, n1 == n2)
	assert.True(t, NodeEquals(n1, n2))

	number1 := NewNumberNode(1.0)
	index1 := NewIndexNode(1.0)
	string1 := NewStringNode("1.0")
	assert.False(t, NodeEquals(number1, index1))
	assert.False(t, NodeEquals(number1, string1))
}

func Test_LessThan(t *testing.T) {
	ten, err := NewNode(10)
	require.Nil(t, err)

	two, err := NewNode(2)
	require.Nil(t, err)

	assert.True(t, two.LessThan(ten))
	assert.True(t, NodeLessThan(two, ten))
}

func Test_unary_function_should_be_a_node(t *testing.T) {
	s := NewNodeSet()
	s.Add(SquareFunctionNode)
	assert.True(t, s.Contains(SquareFunctionNode))
}

func Test_IndexNode_LessThan(t *testing.T) {
	two := NewIndexNode(2)
	one := NewIndexNode(1)
	assert.True(t, one.LessThan(two))
}

func Test_LessThan_Mixed(t *testing.T) {
	one := NewIndexNode(1)
	z := NewStringNode("z")
	assert.True(t, one.LessThan(z))
}
