package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTriplesFromMap(t *testing.T) {
	triples := NewTriples()
	_, err := triples.NewTriplesFromMap(map[string]interface{}{
		"first": "marc",
		"last":  "von Holzen",
		"age":   50,
	})
	assert.Nil(t, err)
	assert.Len(t, triples.TripleSet, 3)
}

func Test_Contains(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("d", "e", 2)
	triple, _ := NewTriple("a", "b", 1)

	assert.True(t, tpls.Contains(triple))
}

// func Test_Contains_unary(t *testing.T) {
// 	tpls := NewTriples()
// 	tpls.AddTriple("a", "b", 1)
// 	tpls.AddTriple("d", "e", 2)
// 	triple, _ := NewTriple("a", NewNodeMatchAny(), 1)

// 	assert.True(t, tpls.Contains(triple))
// }
