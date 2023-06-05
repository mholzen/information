package triples

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unary_function_should_be_a_node(t *testing.T) {
	s := NewNodeSet()
	s.Add(SquareNode)
	assert.True(t, s.Contains(SquareNode))
}

func Test_compute_square(t *testing.T) {
	triples := NewTriples()
	x := NewAnonymousNode()
	triples.NewTriple(x, SquareNode, NewNumberNode(2))
	triples.Compute()
	assert.Len(t, triples.TripleSet, 2)
	log.Printf("triples:\n%s", triples.String())

	assert.Contains(t, triples.TripleSet, Triple{x, NewStringNode(SquareNode.String()), NewNumberNode(4)}.String())

	newTriple := Triple{x, NewStringNode(SquareNode.String()), NewNumberNode(4)}
	assert.True(t, triples.Contains(newTriple))
}

func Test_compute_search(t *testing.T) {
	triples := NewTriples()
	x := NewAnonymousNode()
	triples.NewTriple(x, NewStringNode("first"), NewStringNode("marc"))
	triples.NewTriple(NewStringNode("x?"), NewStringNode("first"), NewStringNode("marc"))
	triples.Compute()
	assert.Len(t, triples.TripleSet, 3)
	assert.Contains(t, triples.TripleSet, Triple{NewStringNode("x"), NewStringNode("equals"), x})
}

// func Test_Compute_Reducer(t *testing.T) {
// 	triples := NewTriples()
// 	_, err := triples.NewTriplesFromSlice([]interface{}{
// 		"x", "sum", "1", "2",
// 	})
// 	assert.Nil(t, err)
// 	assert.Len(t, triples.TripleSet, 2)
// 	triples.Compute()
// 	assert.Len(t, triples.TripleSet, 3)
// 	assert.Contains(t, triples.TripleSet, Triple{NewStringNode("x"), NewStringNode("equals"), NewStringNode("3")})
// }
