package transforms

import (
	"log"
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_Traverse(t *testing.T) {
	graph := triples.NewTriples()
	x := triples.NewAnonymousNode()
	y := triples.NewAnonymousNode()
	z := triples.NewAnonymousNode()
	graph.AddTriple(x, "next", y)
	graph.AddTriple(y, "next", z)

	res := triples.NewTriples()
	tr := NewTraverse(x, NewPredicateTripleMatch(triples.NewStringNode("next")), res)
	err := tr(graph)
	require.Nil(t, err)
	log.Printf("res: %s", res)
	require.Equal(t, 8, len(res.TripleSet))
}

func Test_NodeTraverse(t *testing.T) {
	graph := triples.NewTriples()
	x := triples.NewAnonymousNode()
	y := triples.NewAnonymousNode()
	z := triples.NewAnonymousNode()
	graph.AddTriple(x, "next", y)
	graph.AddTriple(y, "next", z)

	next := func(node triples.Node) triples.NodeList {
		return graph.GetTripleListForSubject(node).GetObjects()
	}

	tail := NewNodeTraverse(x, next)
	res, err := graph.Map(tail)
	require.Nil(t, err)
	require.Equal(t, 3, len(res.TripleSet))

	reverse := func(node triples.Node) triples.NodeList {
		return graph.GetTripleListForObject(node).GetSubjects()
	}
	rootMapper := NewNodeTraverse(z, reverse)
	res, err = graph.Map(rootMapper)
	require.Nil(t, err)
	require.Equal(t, 3, len(res.TripleSet))

	require.Equal(t, res.GetTripleList()[2].Object, x) // Root is the last object

}
