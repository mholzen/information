package triples

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_filter(t *testing.T) {
	data, err := DecodeJson(`{"first":"marc","last":"von Holzen"}`)
	assert.Nil(t, err)

	src := NewTriples()
	err = src.Transform(NewParser(data))
	assert.Nil(t, err)

	fn := NewStringNodeMatch(".*")
	res := NewTriples()
	f := NewNodeFilter(res, fn)
	err = src.Transform(f)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 2)
}

func Test_predicate_filter(t *testing.T) {
	var top *Node
	data, err := NewJsonParser(`{"first":"marc","last":"von Holzen"}`, top)
	assert.Nil(t, err)

	src := NewTriples()
	err = src.Transform(data)
	assert.Nil(t, err)

	res := NewTriples()
	f := NewTripleFilter(res, NewPredicateTripleMatch("first"))
	err = src.Transform(f)
	assert.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)
}

func Test_traverse(t *testing.T) {
	var top Node = NewAnonymousNode()
	tm, err := NewJsonParser(`{"first":["marc","marco"],"last":"von Holzen"}`, &top)
	assert.Nil(t, err)
	assert.NotNil(t, top)

	src := NewTriples()
	err = src.Transform(tm)
	assert.Nil(t, err)
	assert.Len(t, src.TripleSet, 4)
	assert.Len(t, src.GetTriplesForSubject(top), 2)

	res := NewTriples()
	dest := NewAnonymousNode()
	err = src.Transform(NewTraverse(top, AlwaysTripleMatch, dest, res))
	assert.Nil(t, err)
	assert.Len(t, res.GetTriplesForSubject(dest), 4)

	res2 := NewTriples()
	dest2 := NewAnonymousNode()

	// match := NewObjectTripleMatch(NewStringNodeMatch(".*"))

	objectMapper := NewTripleObjectTransformer(dest2, res2)
	err = res.Transform(NewMap(dest, objectMapper, res2))
	assert.Nil(t, err)

	answer := res2.GetTriplesForSubject(dest2)
	log.Printf("answer: %s", answer)
	assert.Len(t, answer, 4)
}
