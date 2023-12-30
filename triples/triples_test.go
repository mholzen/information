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

func Test_GetReference(t *testing.T) {
	tpls := NewTriples()
	tr1, _ := tpls.AddTriple("a", "b", 1)
	tr2, _ := tpls.AddTriple("c", "d", 2)
	n := tpls.AddTripleReference(tr1)

	nodes := tpls.GetTripleReferences(tr1)
	require.Len(t, nodes, 1)
	assert.Contains(t, nodes, n.String())
	require.True(t, nodes.Contains(n))

	nodes = tpls.GetTripleReferences(tr2)
	require.Len(t, nodes, 0)
}

func Test_AddTripleNodes(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTripleNodes("a", "b", 1, "c", "d", 2)
	assert.Len(t, tpls.TripleSet, 2)

	err := tpls.AddTripleNodes("a", "b", 1, "c", "d")
	require.NotNil(t, err)

	err = tpls.AddTripleNodes("a", "b", tpls)
	require.NotNil(t, err)
}
