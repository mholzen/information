package triples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Cartesian(t *testing.T) {
	a := NewTriples()
	a.AddTriple("a", "x", 1)
	a.AddTriple("a", "x", 2)
	b := NewTriples()
	b.AddTriple("b", "x", 3)
	b.AddTriple("b", "x", 4)
	b.AddTriple("b", "x", 5)
	c := NewTriples()
	c.AddTriple("c", "x", 6)
	c.AddTriple("c", "x", 7)

	sets := TriplesList{a, b, c}

	res := sets.Cartesian()
	require.Len(t, res, 12)
}

func Test_Cartesian_empty_set(t *testing.T) {
	a := NewTriples()
	b := NewTriples()
	b.AddTriple("a", "b", 1)

	sets := TriplesList{a, b}

	res := sets.Cartesian()
	require.Len(t, res, 1)
}
