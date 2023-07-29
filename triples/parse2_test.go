package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse2(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)

	assert.Len(t, src.TripleSet, 2)
}

func Test_parse_with_result(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.NotNil(t, tm.Result)

	assert.Len(t, src.TripleSet, 2)
}
