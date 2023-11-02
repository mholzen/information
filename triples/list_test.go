package triples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Sort(t *testing.T) {
	triples := TripleList{
		NewTripleFromNodes(NewStringNode("a"), NewIndexNode(2), NewStringNode("c")),
		NewTripleFromNodes(NewStringNode("a"), NewIndexNode(1), NewStringNode("b")),
		NewTripleFromNodes(NewStringNode("a"), NewIndexNode(1), NewStringNode("a")),
	}

	sorted := triples.Sort()
	require.Equal(t, "a", sorted[0].Object.String())
}

func Test_SortObjects(t *testing.T) {
	triples := TripleList{
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("a"), NewIndexNode(4)),
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("a"), NewIndexNode(1)),
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("a"), NewIndexNode(2)),
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("a"), NewIndexNode(3)),
	}

	sorted := triples.Sort()
	require.Equal(t, "1", sorted[0].Object.String())
}
