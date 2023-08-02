package triples

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_filter(t *testing.T) {
	data, err := DecodeJson(`{"first":"marc","last":"von Holzen"}`)
	assert.Nil(t, err)

	src := NewTriples()
	err = src.Transform(NewParser(data))
	assert.Nil(t, err)

	fn := NewStringNodeMatch("von.*")
	res := NewTriples()
	err = src.Transform(NewObjectFilter(res, fn))
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 1)
}

func Test_predicate_filter(t *testing.T) {
	data := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(data.Transformer)
	assert.Nil(t, err)

	res := NewTriples()
	first, _ := NewNode("first") // Can avoid using a function that returns an error?
	f := NewTripleFilter(res, NewPredicateTripleMatch(first))
	err = src.Transform(f)
	assert.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)
}

func Test_traverse(t *testing.T) {
	tm := NewJsonParser(`{"first":["marc","marco"],"last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)
	assert.Len(t, src.TripleSet, 4)
	assert.Len(t, src.GetTriplesForSubject(*tm.Result), 2)

	res := NewTriples()
	dest := NewAnonymousNode()
	err = src.Transform(NewTraverse(*tm.Result, AlwaysTripleMatch, dest, res))
	assert.Nil(t, err)
	assert.Len(t, res.GetTriplesForSubject(dest), 4)

	res2 := NewTriples()
	dest2 := NewAnonymousNode()
	objectMapper := NewTripleObjectTransformer(dest2, res2)
	err = res.Transform(NewMap(dest, objectMapper, res2))
	assert.Nil(t, err)

	answer := res2.GetTriplesForSubject(dest2)
	assert.Len(t, answer, 4)
}

func Test_traverse_file(t *testing.T) {
	var top Node = NewAnonymousNode()
	tm := NewFileJsonParser("../data/verbs.jsonc")

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	assert.Nil(t, err)

	dest := NewAnonymousNode()
	err = src.Transform(NewTraverse(top, AlwaysTripleMatch, dest, src))
	assert.Nil(t, err)

	dest2 := NewAnonymousNode()
	objectMapper := NewTripleObjectTransformer(dest2, src)
	err = src.Transform(NewMap(dest, objectMapper, src))
	assert.Nil(t, err)

	res := NewTriples()
	err = src.Transform(NewFlatMap(dest2, GetStringObjectMapper, res))
	assert.Nil(t, err)

	answer := res.GetTripleList().GetObjectStrings()
	sort.Strings(answer)

	assert.Greater(t, len(src.TripleSet), 100)
}
