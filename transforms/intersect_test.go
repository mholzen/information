package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_Intersect(t *testing.T) {
	tpls, err := triples.NewTriplesFromAny(
		"a", "b", 1,
		"a", "b", 2,
		"a", "b", 3,
	)
	require.Nil(t, err)

	toRemove, err := triples.NewTriplesFromAny(
		"a", "b", 2,
	)
	require.Nil(t, err)

	intersect := NewIntersectMapper(toRemove)
	result, err := intersect(tpls)
	require.Nil(t, err)

	require.Contains(t, result.TripleSet, "(a, b, 1)")
	require.NotContains(t, result.TripleSet, "(a, b, 2)")
	require.Contains(t, result.TripleSet, "(a, b, 3)")
}
