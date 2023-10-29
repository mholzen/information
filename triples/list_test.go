package triples

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_Sort(t *testing.T) {
	triples := TripleList{
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("b"), NewStringNode("d")),
		NewTripleFromNodes(NewStringNode("a"), NewStringNode("b"), NewStringNode("c")),
	}

	sorted := triples.Sort()
	require.Equal(t, sorted[0].Object.String(), "c")
}
