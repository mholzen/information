package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Computation(t *testing.T) {
	NewComputation(NewVariableNode(), TypeFunctionNode, NewStringNode("triples.AnonymousNode"))

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
