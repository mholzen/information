package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_filter(t *testing.T) {
	data, err := DecodeJson(`{"first":"marc","last":"von Holzen"}`)
	require.Nil(t, err)

	src := NewTriples()
	err = src.Transform(NewParser(data))
	require.Nil(t, err)

	fn := NewStringNodeMatch("von.*")
	res := NewTriples()
	err = src.Transform(NewObjectFilter(res, fn))
	require.Nil(t, err)
	assert.Len(t, res.TripleSet, 1)
}

func Test_predicate_filter(t *testing.T) {
	data := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(data.Transformer)
	require.Nil(t, err)

	res := NewTriples()
	first, _ := NewNode("first") // Can avoid using a function that returns an error?
	f := NewFilterTransformer(res, NewPredicateTripleMatch(first))
	err = src.Transform(f)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)
}

func Test_traverse(t *testing.T) {
	tm := NewJsonParser(`{"first":["marc","marco"],"last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.Len(t, src.TripleSet, 4)
	assert.Len(t, src.GetTripleListForSubject(*tm.Result), 2)

	res := NewTriples()
	err = src.Transform(NewTraverse(*tm.Result, AlwaysTripleMatch, res))
	require.Nil(t, err)
	assert.Len(t, res.TripleSet, 4 /*nodes*/ +4*3 /*references per node*/)

	references, err := GetReferences(res)
	require.Nil(t, err)
	assert.Len(t, references.TripleSet, 4)
}

func Test_traverse_file(t *testing.T) {
	tm := NewFileJsonParser("../../data/verbs.jsonc")

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.Greater(t, len(src.TripleSet), 100)

	dest := NewTriples()
	err = src.Transform(NewTraverse(*tm.Result, AlwaysTripleMatch, dest))
	require.Nil(t, err)

	assert.Greater(t, len(dest.TripleSet), 10)
}
