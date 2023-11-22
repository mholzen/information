package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCompute(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", TypeFunctionNode, "foo")

	err := NewCompute()(tpls)
	require.Nil(t, err)

	assert.Len(t, tpls.TripleSet, 2)
	shouldContain, _ := NewTriple("a", NewStringNode(TypeFunctionNode.String()), "triples.StringNode")
	assert.Contains(t, tpls.GetTripleList(), shouldContain)
	assert.True(t, tpls.Contains(shouldContain))
}

func TestComputeSquare(t *testing.T) {
	triples := NewTriples()
	triples.NewTripleFromNodes(NewNumberNode(2), NewStringNode("square"), NewVariableNode())

	computer := NewComputeWithDefinitions(GetDefinitions())
	err := computer(triples)
	require.Nil(t, err)

	assert.Len(t, triples.TripleSet, 2)

	shouldContain, err := NewTriple(NewNumberNode(2), "square", NewNumberNode(4))
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func TestComputeType(t *testing.T) {
	triples := NewTriples()
	x := NewAnonymousNode()
	triples.NewTripleFromNodes(x, NewStringNode("type"), NewVariableNode())
	triples.NewTripleFromNodes(NewStringNode("foo"), NewStringNode("type"), NewVariableNode())
	triples.NewTripleFromNodes(NewIndexNode(1), NewStringNode("type"), NewVariableNode())
	triples.NewTripleFromNodes(NewNumberNode(3.14), NewStringNode("type"), NewVariableNode())

	computer := NewComputeWithDefinitions(GetDefinitions())
	err := computer(triples)
	require.Nil(t, err)

	// log.Printf("triples are:\n%v", triples)
	shouldContain, err := NewTriple(x, "type", "triples.AnonymousNode")
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func TestComputeWithDefinitions(t *testing.T) {
	triples := NewTriples()
	triples.NewTripleFromNodes(NewNumberNode(2), NewStringNode("square"), NewVariableNode())

	definitions := NewTriples()
	definitions.NewTripleFromNodes(SquareFunctionNode, ComputeNode, NewStringNode("square"))

	tr := NewComputeWithDefinitions(definitions)
	err := tr(triples)
	require.Nil(t, err)

	assert.Len(t, triples.TripleSet, 2)

	shouldContain, err := NewTriple(NewNumberNode(2), "square", NewNumberNode(4))
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func Test_compute_wanted_triples(t *testing.T) {
	// a triple expressing desired result
	// x, square, 4

	// a triple expressing the label to store the result of the function
	// square, computed-by, square-f

}

func TestComputeTripleTransformer(t *testing.T) {

}
