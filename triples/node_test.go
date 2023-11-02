package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Compare(t *testing.T) {
	a1, _ := NewNode("a")
	a2, _ := NewNode("a")
	assert.True(t, a1 == a2)

	n1, _ := NewNode(1)
	n2, _ := NewNode(1)
	assert.True(t, n1 == n2)
}

func Test_LessThan(t *testing.T) {
	ten, err := NewNode(10)
	require.Nil(t, err)
	two, err := NewNode(2)
	require.Nil(t, err)
	assert.True(t, two.LessThan(ten))
}

func Test_unary_function_should_be_a_node(t *testing.T) {
	s := NewNodeSet()
	s.Add(SquareNode)
	assert.True(t, s.Contains(SquareNode))
}

func Test_IndexNode_LessThan(t *testing.T) {
	two := NewIndexNode(2)
	one := NewIndexNode(1)
	assert.True(t, one.LessThan(two))
}

func Test_LessThan_Mixed(t *testing.T) {
	one := NewIndexNode(1)
	a := NewStringNode("z")
	assert.True(t, one.LessThan(a))
}
