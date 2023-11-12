package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_CompareTriple(t *testing.T) {
	t1, _ := NewTriple("a", "b", 1)
	t2, _ := NewTriple("a", "b", 1)
	t3, _ := NewTriple("a", "c", 1)

	assert.Equal(t, t1, t2)
	assert.NotEqual(t, t1, t3)
}

func Test_NewTriplesFromMap(t *testing.T) {
	triples := NewTriples()
	_, err := triples.NewTriplesFromMap(map[string]interface{}{
		"first": "marc",
		"last":  "von Holzen",
		"age":   50,
	})
	require.Nil(t, err)
	assert.Len(t, triples.TripleSet, 3)
}

func Test_Contains(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("d", "e", 2)
	triple, _ := NewTriple("a", "b", 1)

	assert.True(t, tpls.Contains(triple))
	assert.Contains(t, tpls.TripleSet, triple.String())

	t1 := NewTriples()
	shouldContain, err := NewTriple(NewNumberNode(2), "square", NewNumberNode(4))
	require.Nil(t, err)
	t1.Add(shouldContain)
	assert.Contains(t, t1.TripleSet, shouldContain.String())

}

// func Test_Contains_unary(t *testing.T) {
// 	tpls := NewTriples()
// 	tpls.AddTriple("a", "b", 1)
// 	tpls.AddTriple("d", "e", 2)
// 	triple, _ := NewTriple("a", NewNodeMatchAny(), 1)

// 	assert.True(t, tpls.Contains(triple))
// }
