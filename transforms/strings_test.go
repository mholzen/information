package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ObjectsToStrings(t *testing.T) {
	res := triples.NewTriples()
	_, err := res.AddTriple(triples.NewStringNode("a"), triples.NewIndexNode(1), triples.NewIndexNode(2))
	require.NoError(t, err)
	res2, _ := ObjectsToStrings(res)

	o := res2.GetObjects().GetNodeList()[0]

	assert.IsType(t, triples.StringNode{}, o)
}
