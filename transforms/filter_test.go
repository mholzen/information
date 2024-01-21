package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_TripleMatch(t *testing.T) {
	src := triples.NewTriples()
	triple, _ := src.AddTripleFromAny("a", "b", 1)
	require.True(t, NewSubjectTripleMatch(triples.NewStringNode("a"))(triple))
	require.True(t, NewPredicateTripleMatch(triples.NewStringNode("b"))(triple))
	require.True(t, NewObjectTripleMatch(triples.NewIndexNode(1))(triple))
}

func Test_FilterForIndexNodes(t *testing.T) {
	src := triples.NewTriples()
	src.AddTripleFromAny("a", "b", 1)
	src.AddTripleFromAny("a", "b", 2)
	src.AddTripleFromAny("b", 2, "c")

	filter := triples.NewTriples()
	container := triples.NewAnonymousNode()
	// filter.AddTriple(container, triples.Subject, triples.NewNodeMatchAny())
	_, err := filter.AddTripleFromAny(container, triples.Predicate, triples.NodeMatchAnyIndex)
	require.Nil(t, err)
	// filter.AddTriple(container, triples.Subject, triples.NewNodeMatchAny())

	// res, err := src.Map(NewFilterMapperFromTriples(filter))
	tripleMatch, err := NewTripleMatchFromTriples(filter)
	require.Nil(t, err)

	res, err := src.Map(Filter(tripleMatch))
	require.Nil(t, err)
	require.Len(t, res.TripleSet, 1)
}

func Test_FilterForSubjectPredicate(t *testing.T) {
	src := triples.NewTriples()
	subject := triples.NewAnonymousNode()
	src.AddTripleFromAny(subject, "first", "marc")

	filter := triples.NewTriples()
	filter.AddTripleFromAny(triples.NewAnonymousNode(), triples.Subject, subject)
	filter.AddTripleFromAny(triples.NewAnonymousNode(), triples.Predicate, "first")

	m, err := NewTripleMatchFromTriples(filter)
	require.Nil(t, err)

	res, err := src.Map(Filter(m))
	require.Nil(t, err)
	require.Len(t, res.TripleSet, 1)
}
