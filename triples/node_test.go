package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_node_compare(t *testing.T) {
	a1, _ := NewNode("a")
	a2, _ := NewNode("a")
	assert.True(t, a1 == a2)

	n1, _ := NewNode(1)
	n2, _ := NewNode(1)
	assert.True(t, n1 == n2)
}
