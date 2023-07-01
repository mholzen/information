package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse2(t *testing.T) {
	var top *Node
	tm, err := NewJsonParser(`{"first":"marc","last":"von Holzen"}`, top)
	assert.Nil(t, err)

	src := NewTriples()
	err = src.Transform(tm)
	assert.Nil(t, err)

	assert.Len(t, src.TripleSet, 2)
}
