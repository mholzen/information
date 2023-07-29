package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.NotNil(t, tm.Result)

	assert.Len(t, src.TripleSet, 2)
}

func Test_triples_map2(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}]`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.Len(t, src.TripleSet, 2)
}

func Test_triples_array_object(t *testing.T) {
	tm := NewJsonParser(`{"names":["marc","Marc", "Marco"]}`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.Len(t, src.TripleSet, 4)
}

func Test_slice_as_object(t *testing.T) {
	tm := NewJsonParser(`["root", "contains", ["marc", "is", "alive"]]`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.Len(t, src.TripleSet, 6)
}

func Test_csv_parse(t *testing.T) {
	tm := NewCsvParser("a,b,c\nd,e,f")

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)

	assert.Len(t, src.TripleSet, 8)
}
