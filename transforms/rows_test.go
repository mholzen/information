package transforms

import (
	"log"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_RowQuery(t *testing.T) {
	query := RowQuery()
	source, err := NewTriplesFromStrings(
		"_ 1 _",
		"_ 2 _",
		"_ 3 _",
		"_ 4 1",
		"a 5 _",
	)
	require.Nil(t, err)

	solutions, err := query.SearchForSolutions(source)
	require.Nil(t, err)

	require.Len(t, solutions, 3)
}

func Test_MatrixQuery(t *testing.T) {
	// TODO: fail if NewTriplesFromStrings
	source, err := NewNamedTriples(
		"_root 1 _row1",
		"_root 2 _row2",
		"_row1 1 _cell11",
		"_row1 2 _cell12",
		"_row2 1 _cell21",
		"_row2 2 _cell22",
		"_cell11 a 11",
		"_cell11 b 11",
		"_cell12 a 12",
		"_cell12 b 12",
		"_cell21 a 21",
		"_cell21 b 21",
		"_cell22 a 22",
		"_cell22 b 22",
	)
	require.Nil(t, err)

	query := MatrixQuery()
	matchesMap, err := query.SearchForMatches(source)
	require.Nil(t, err)
	require.Len(t, matchesMap, 2)

	require.Len(t, matchesMap[query.Matching.GetTripleList()[0]].TripleSet, source.Length())

	solutions, err := query.SearchForSolutions(source)
	require.Nil(t, err)

	log.Printf("%s", solutions.GetAllTriples())
	require.Len(t, solutions, 4)

	selected, err := solutions.GetSelectTriples(query.Selected)
	require.Nil(t, err)

	log.Printf("select: %s", selected)
	require.Len(t, selected.TripleSet, 2)
}
