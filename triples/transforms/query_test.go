package transforms

import (
	"testing"

	"github.com/labstack/gommon/log"
	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewQueryMapperWithMatchesAndJoins(t *testing.T) {
	t.Skip()
	tpls := NewTriples()
	marc1 := NewAnonymousNode()
	tpls.AddTriple(marc1, "first", "john")
	tpls.AddTriple(marc1, "age", 52)

	marc2 := NewAnonymousNode()
	tpls.AddTriple(marc2, "first", "marc")
	tpls.AddTriple(marc2, "age", 16)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, "first", NodeMatchAnyString) // VariableNode matches anything but are used to join
	query.AddTriple(x, "age", 52)

	mapper, err := NewQueryMapper(query)
	require.Nil(t, err)
	res, err := tpls.Map(mapper)
	require.Nil(t, err)
	log.Debug(res)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2) // TODO: should join on both query triples -- iterate over anonymous nodes in the query
	assert.True(t, refs.ContainsTriple(marc1, "age", "52"))
}

func Test_Computation(t *testing.T) {
	NewComputation(NewVariableNode(), TypeFunctionNode, NewStringNode("triples.AnonymousNode"))

}

func Test_NewQueryWithSimpleComputations(t *testing.T) {
	t.Skip()
	tpls := NewTriples()
	tpls.AddTriple("marc", "name", "Marc")
	marc := NewAnonymousNode()
	tpls.AddTriple(marc, "name", "Marc")
	tpls.AddTriple(NewAnonymousNode(), "name", "John")

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, "name", "Marc")
	query.AddTriple(x, UnaryFunctionNode(TypeFunctionNode), "triples.AnonymousNode")

	q, err := NewQueryFromTriples(query)
	require.Nil(t, err)
	solutions, err := q.Apply(tpls)
	require.Nil(t, err)
	require.Len(t, solutions, 1)
}

func Test_NewQueryWithComputationsComplex(t *testing.T) {
	t.Skip()
	tpls := NewTriples()
	tpls.AddTriple("marc", "age", 10)
	tpls.AddTriple("john", "age", 4)

	query := NewTriples()
	squaredAge := NewAnonymousNode()
	query.AddTriple(squaredAge, "function", SquareFunctionNode)

	query.AddTriple(NewVariableNode(), squaredAge, 16)

	mapper, err := NewQueryMapper(query)
	require.Nil(t, err)
	res, err := tpls.Map(mapper)
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2)
	assert.True(t, refs.ContainsTriple("a", "b", 1))
	assert.True(t, refs.ContainsTriple("c", "d", 1))
}

func people() *Triples {
	tpls := NewTriples()
	marc, _ := tpls.AddTriple(NewAnonymousNode(), "name", "Marc")
	tpls.AddTriple(marc.Subject, "age", 50)
	john, _ := tpls.AddTriple(NewAnonymousNode(), "name", "John")
	tpls.AddTriple(john.Subject, "age", 24)
	marry, _ := tpls.AddTriple(NewAnonymousNode(), "name", "Marry")
	tpls.AddTriple(marry.Subject, "age", 32)

	return tpls
}
func Test_QueryTripleMatcher_Simple(t *testing.T) {
	query, _ := NewTriple(NewVariableNode(), "name", NodeMatchAnyString)

	res, err := people().Map(NewTripleQueryMatchMapper(query))
	require.Nil(t, err)
	require.Len(t, res.TripleSet, 3)
}

func Test_QueryTripleMatcher_Compute(t *testing.T) {
	t.Skip()
	query, _ := NewTriple(NewVariableNode(), LengthFunction, 4)
	mapper := NewTripleQueryMatchMapper(query)
	// TODO: compute doesn't apply to a triple (so it's not a mapper), it applies to nodes, probably to solutions

	tpls := NewTriples()
	tpls.AddTriple(NewAnonymousNode(), "name", "Marc")

	res, err := mapper(tpls)
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 1)
	assert.True(t, refs.ContainsTriple("a", "b", 1))
}
