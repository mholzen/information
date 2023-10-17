package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_SortedSet(t *testing.T) {
	set := NewNodeSet()
	set.Add(NewStringNode("b"))
	set.Add(NewStringNode("a"))
	set.Add(NewStringNode("c"))
	// TODO: should sort nodes from different types

	nodes := set.GetSortedNodeList()
	assert.Equal(t, nodes[0], NewStringNode("a"))
	assert.Equal(t, nodes[1], NewStringNode("b"))
	assert.Equal(t, nodes[2], NewStringNode("c"))
}
