package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_Cartesian(t *testing.T) {
	a := tr.NewTriples()
	a.AddTriple("a", "x", 1)
	a.AddTriple("a", "x", 2)
	b := tr.NewTriples()
	b.AddTriple("b", "x", 3)
	b.AddTriple("b", "x", 4)
	b.AddTriple("b", "x", 5)
	c := tr.NewTriples()
	c.AddTriple("c", "x", 6)
	c.AddTriple("c", "x", 7)

	sets := []*tr.Triples{a, b, c}

	res := Cartesian(sets)
	require.Len(t, res, 12)
}

func Test_Cartesian_empty_set(t *testing.T) {
	a := tr.NewTriples()
	b := tr.NewTriples()
	b.AddTriple("a", "b", 1)

	sets := []*tr.Triples{a, b}

	res := Cartesian(sets)
	require.Len(t, res, 1)
}
