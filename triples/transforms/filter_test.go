package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_FilterForIndexNodes(t *testing.T) {
	src := triples.NewTriples()
	src.AddTriple("a", "b", 1)
	src.AddTriple("a", "b", 2)
	src.AddTriple("b", 2, "c")

	filter := triples.NewTriples()
	container := triples.NewAnonymousNode()
	// filter.AddTriple(container, triples.Subject, triples.NewNodeMatchAny())
	filter.AddTriple(container, triples.Predicate, triples.NewNodeMatchAnyIndex())
	// filter.AddTriple(container, triples.Subject, triples.NewNodeMatchAny())

	res, err := src.Map(NewFilterMapperFromTriples(filter))
	require.Nil(t, err)
	require.Len(t, res.TripleSet, 1)
}

func Test_FilterForSubjectPredicate(t *testing.T) {
	src := triples.NewTriples()
	subject := triples.NewAnonymousNode()
	src.AddTriple(subject, "first", "marc")

	filter := triples.NewTriples()
	filter.AddTriple(triples.NewAnonymousNode(), triples.Subject, subject)
	filter.AddTriple(triples.NewAnonymousNode(), triples.Predicate, "first")

	m, err := NewTripleMatchFromTriples(filter)
	require.Nil(t, err)

	res, err := src.Map(Filter(m))
	require.Nil(t, err)
	require.Len(t, res.TripleSet, 1)
}
