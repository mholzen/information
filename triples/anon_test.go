package triples

import (
	"testing"

	"github.com/subchen/go-stack/assert"
)

func Test_NewAnonNode(t *testing.T) {
	anonNode1 := NewAnonymousNode()
	anonNode2 := NewAnonymousNode()
	anonNode3 := NewAnonymousNode()

	assert.True(t, anonNode1.LessThan(anonNode2))
	assert.True(t, anonNode1.LessThan(anonNode3))
	assert.True(t, anonNode2.LessThan(anonNode3))
}
