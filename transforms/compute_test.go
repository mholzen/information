package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewCompute(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTripleFromAny("a", TypeFunctionNode, "foo")

	err := NewCompute()(tpls)
	require.Nil(t, err)

	assert.Len(t, tpls.TripleSet, 2)
	var f UnaryFunctionNode = TypeFunctionNode
	shouldContain, _ := NewTripleFromAny("a", f.String(), "triples.StringNode")
	assert.Contains(t, tpls.GetTripleList(), shouldContain)
	assert.True(t, tpls.Contains(shouldContain))
}

func TestComputeSquare(t *testing.T) {
	triples := NewTriples()
	triples.AddTriple(NewNumberNode(2), NewStringNode("square"), NewVariableNode())

	computer := NewComputeWithDefinitions(GetDefinitions())
	err := computer(triples)
	require.Nil(t, err)

	assert.Len(t, triples.TripleSet, 2)

	shouldContain, err := NewTripleFromAny(NewNumberNode(2), "square", NewNumberNode(4))
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func TestComputeType(t *testing.T) {
	triples := NewTriples()
	x := NewAnonymousNode()
	triples.AddTriple(x, NewStringNode("type"), NewVariableNode())
	triples.AddTriple(NewStringNode("foo"), NewStringNode("type"), NewVariableNode())
	triples.AddTriple(NewIndexNode(1), NewStringNode("type"), NewVariableNode())
	triples.AddTriple(NewNumberNode(3.14), NewStringNode("type"), NewVariableNode())

	computer := NewComputeWithDefinitions(GetDefinitions())
	err := computer(triples)
	require.Nil(t, err)

	// log.Printf("triples are:\n%v", triples)
	shouldContain, err := NewTripleFromAny(x, "type", "triples.AnonymousNode")
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func TestComputeWithDefinitions(t *testing.T) {
	triples := NewTriples()
	triples.AddTriple(NewNumberNode(2), NewStringNode("square"), NewVariableNode())

	definitions := NewTriples()
	definitions.AddTriple(SquareFunctionNode, ComputeNode, NewStringNode("square"))

	tr := NewComputeWithDefinitions(definitions)
	err := tr(triples)
	require.Nil(t, err)

	assert.Len(t, triples.TripleSet, 2)

	shouldContain, err := NewTripleFromAny(NewNumberNode(2), "square", NewNumberNode(4))
	require.Nil(t, err)
	assert.Contains(t, triples.TripleSet, shouldContain.String())
}

func Test_compute_wanted_triples(t *testing.T) {
	// a triple expressing desired result
	// x, square, 4

	// a triple expressing the label to store the result of the function
	// square, computed-by, square-f

}

func Test_TripleMapper(t *testing.T) {
	SubjectTypeMapper, err := NewPositionFunctionMapper(SubjectPosition, TypeFunctionNode, NewStringNode("type"))
	require.Nil(t, err)

	triples := NewTriples()
	a := NewAnonymousNode()
	triples.AddTripleFromAny(a, "v", 1)
	triples.AddTripleFromAny("b", "v", 2)
	triples.AddTripleFromAny(2, "v", 3)

	res, err := triples.MapTriples(SubjectTypeMapper)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 3)
	assert.Contains(t, res.String(), "type, triples.AnonymousNode)")
	assert.Contains(t, res.String(), "(b, type, triples.StringNode)")
	assert.Contains(t, res.String(), "(2, type, triples.IndexNode)")
}

func Test_SubjectFunctionMapperFromTriples(t *testing.T) {

	query := NewTriple(NewVariableNode(), UnaryFunctionNode(TypeFunctionNode), Str("triples.AnonymousNode"))
	generator, err := NewSubjectFunctionGeneratorFromTriples(query)
	require.Nil(t, err)

	tpls := NewTriples()
	tpls.AddTripleFromAny(NewAnonymousNode(), "is", "anonymous")
	tpls.AddTripleFromAny("text", "is", "string")

	res, err := tpls.FlatMap(generator)
	require.Nil(t, err)

	require.Len(t, res.TripleSet, 1)
	assert.Contains(t, res.String(), "is, anonymous)")
}
