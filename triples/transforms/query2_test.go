package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewQuery2(t *testing.T) {
	data, err := NewTriplesFromStrings(
		"a b 1",
		"a b 2",
		"c d 1",
		"c d 2",
	)
	require.Nil(t, err)

	queryTriples, err := NewTriplesFromStrings(
		"a b 2",
		"c d 2",
	)
	require.Nil(t, err)
	query := NewQuery2(queryTriples, Computations{})

	solutions, err := query.Apply(data)
	require.Nil(t, err)

	assert.Len(t, solutions, 1)
	solution := solutions[0]
	assert.Len(t, solution.SolutionMap, 2)

	assert.True(t, solutions.GetAllTriples().ContainsTriple("a", "b", 2))
	assert.True(t, solutions.GetAllTriples().ContainsTriple("c", "d", 2))
	assert.False(t, solutions.GetAllTriples().ContainsTriple("a", "b", 1))
}

func Test_NewQuery2_Variables(t *testing.T) {
	data, err := NewTriplesFromStrings(
		"a b 1",
		"a b 2",
		"c d 1",
		"c d 2",
	)
	require.Nil(t, err)

	queryTriples, err := NewTriplesFromStrings(
		"? b 2",
		"c ? 2",
	)
	require.Nil(t, err)
	query := NewQuery2(queryTriples, Computations{})

	solutions, err := query.Apply(data)
	require.Nil(t, err)

	require.Len(t, solutions, 1)
	solution := solutions[0]
	assert.Len(t, solution.SolutionMap, 2)
}

func Test_NewQuery2_Variables_Joins(t *testing.T) {
	data, err := NewTriplesFromStrings(
		"a b 1",
		"a b 2",
		"c b 1",
		"c b 3",
	)
	require.Nil(t, err)

	queryTriples, err := NewNamedTriples(
		"?x b 1",
		"?x b 2",
	)
	require.Nil(t, err)

	query := NewQuery2(queryTriples, Computations{})

	solutions, err := query.Apply(data)
	require.Nil(t, err)

	require.Len(t, solutions, 1)
	solution := solutions[0]
	assert.Len(t, solution.SolutionMap, 2)
}

func Test_NewQuery2_Compute(t *testing.T) {
	data, err := NewTriplesFromStrings(
		"a b 1",
		"a b 2",
		"aa b 2",
		"c d 1",
		"c d 2",
	)
	require.Nil(t, err)

	a := Var()
	d := Var()
	queryTriples, err := tr.NewTriplesFromNodes(
		a, "b", 2,
		"c", d, 2,
	)
	require.Nil(t, err)
	computations := NewComputation(a, tr.LengthFunctionNode, tr.NewIndexNode(1))
	query := NewQuery2(queryTriples, NewComputations(computations))

	solutions, err := query.Apply(data)
	require.Nil(t, err)

	require.Len(t, solutions, 1)
	solution := solutions[0]
	assert.Len(t, solution.SolutionMap, 2)
}

func Test_QueryMatches(t *testing.T) {
	queryFirst, err := NewTripleFromString("? first Marc")
	require.Nil(t, err)

	queryAge, err := NewTripleFromString("? age ?")
	require.Nil(t, err)

	queryTriples := tr.NewTriples().AddTripleList(queryFirst, queryAge)

	data, err := NewTriplesFromStrings(
		"_ first Marc",
		"_ age 42",
		"_ first John",
		"_ age 18",
	)
	require.Nil(t, err)

	query := NewQuery2(queryTriples, Computations{})
	matches, err := query.GetMatches(data)
	require.Nil(t, err)

	matchesMap := matches.MatchesMap
	require.Len(t, matchesMap, 2)
	assert.Len(t, matchesMap[queryFirst].TripleSet, 1)
	assert.Len(t, matchesMap[queryAge].TripleSet, 2)

	solutions := query.GetSolutions(matches)
	assert.Len(t, solutions, 2)

	objects := solutions.GetTriples(queryFirst).GetObjects()
	assert.Len(t, objects, 1)
	assert.Contains(t, objects, "Marc")

	objects = solutions.GetTriples(queryAge).GetObjects()
	assert.Len(t, objects, 2)
	assert.Contains(t, objects, "18")
	assert.Contains(t, objects, "42")
}
