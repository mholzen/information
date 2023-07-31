package node_link

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
)

func Test_node_link(t *testing.T) {
	tpls := triples.NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("a", "b", 2)
	tr := NewNodeLinkTransformer()
	err := tpls.Transform(tr.Transformer)
	assert.Nil(t, err)
	assert.NotNil(t, tr.Result)
	assert.Len(t, tr.Result.Nodes, 3)
	assert.Len(t, tr.Result.Links, 2)
}
