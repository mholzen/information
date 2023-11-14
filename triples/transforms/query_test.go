package transforms

import (
	"log"
	"strings"
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VariableNode(t *testing.T) {
	x := NewVariableNode()
	var y Node = x
	_, ok := y.(VariableNode)
	assert.True(t, ok)
}

func Test_VariableList_Traverse(t *testing.T) {
	nodes := NodeList{NewStringNode("a"), NewStringNode("b"), NewStringNode("c")}
	vars := VariableList{NewVariableNode(), NewVariableNode()}
	r := vars.Traverse(nodes)
	assert.Len(t, r, 9)
}

func Test_NewQueryMapper(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("a", "b", 2)
	tpls.AddTriple("c", "d", 1)
	tpls.AddTriple("c", "d", 2)

	query := NewTriples()
	query.AddTriple("a", "b", 1)
	query.AddTriple("c", "d", 1)

	res, err := tpls.Map(NewQueryMapper(query))
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2)
	assert.True(t, refs.ContainsTriple("a", "b", 1))
	assert.True(t, refs.ContainsTriple("c", "d", 1))
}

func Test_NewQueryMapperWithMatches(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("a", "b", 2)
	tpls.AddTriple("c", "d", 1)
	tpls.AddTriple("c", "d", 2)
	tpls.AddTriple("c", "d", "c")

	query := NewTriples()
	query.AddTriple("c", "d", NodeMatchIndex)

	res, err := tpls.Map(NewQueryMapper(query))
	require.Nil(t, err)

	refs := References(res)
	require.Len(t, refs.TripleSet, 2)
	assert.True(t, refs.ContainsTriple("c", "d", 1))
	assert.True(t, refs.ContainsTriple("c", "d", 2))
	assert.False(t, refs.ContainsTriple("c", "d", "c"))
}

func Test_NewQuery_with_multiple_conditions(t *testing.T) {
	t.Skip()
	src, err := NewJsonTriples(`[{"a":1,"b":2},{"b":2},{"c":3}]`)
	require.Nil(t, err)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, NewStringNode("a"), NewFloatNode(1))
	query.AddTriple(x, NewStringNode("b"), NewFloatNode(2))

	res := NewTriples()
	err = src.Transform(NewQueryTransformerWithDefinitions(query, res, GetDefinitions()))
	require.Nil(t, err)
	res, err = res.Map(ReferencesMapper)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 2)
}

// func Test_NewQuery_with_unary_functions(t *testing.T) {
// 	src, err := NewJsonTriples(`[{"a":1},{"b":2},{"c":3}]`)
// 	require.Nil(t, err)

// 	query := NewTriples()
// 	query.AddTriple(NewVariableNode(), NewStringNodeMatch("."), NewNodeMatchAny())

// 	res := NewTriples()
// 	err = src.Transform(NewQuery(query, res))

// 	require.Nil(t, err)
// 	assert.Len(t, res.TripleSet, 3)
// }

// func Test_match_filename_top(t *testing.T) {
// 	query := NewTriples()
// 	x := NewVariableNode()
// 	query.AddTriple(x, "filename", NewStringNodeMatch(".*"))
// 	query.AddTriple(x, NewIndexNode(1), NewNodeMatchAny())

// 	data := NewTriples()
// 	tm := NewFileJsonParser("../data/array.jsonc")
// 	err := data.Transform(tm.Transformer)
// 	require.Nil(t, err)

// 	res := NewTriples()
// 	err = data.Transform(NewQueryMatch(res, query))
// 	require.Nil(t, err)
// 	assert.Len(t, res.TripleSet, 3)

// }

// Need to query with matching node so that I can query for only the nodes that are rows from a CSV file
// - Need to test if triples contain another triple, using matching nodes
func Test_NewQuery_with_matching(t *testing.T) {
	t.Skip()
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("b", "c", 2)
	tpls.AddTriple("e", "f", 3)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, NewNodeMatchAny(), 2)

	res := NewTriples()
	err := NewQueryTransformerWithDefinitions(query, res, GetDefinitions())(tpls)
	require.Nil(t, err)
	assert.Greater(t, len(res.TripleSet), 1) // ContainsTriples modifies the triples

	res, err = res.Map(ReferencesMapper)
	require.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)
	assert.Equal(t, res.String(), "(b, c, 2)\n")

	assert.True(t, res.Contains(Triple{
		Subject:   NewStringNode("b"),
		Predicate: NewStringNode("c"),
		Object:    NewIndexNode(2),
	}))
}

func Test_NewQuery_for_csv(t *testing.T) {
	t.Skip()
	// Load a CSV file
	tm := NewCsvParser(strings.NewReader("a,b,c\nd,e,f"))

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)

	query := NewTriples()
	file := NewVariableNode()
	query.AddTriple(file, "type", "triples.AnonymousNode")
	query.AddTriple(file, "source", "CsvParser")

	// query for lines
	lines := NewTriples()
	err = NewQueryTransformerWithDefinitions(query, lines, GetDefinitions())(src)
	require.Nil(t, err)

	// Count solutions
	log.Printf("lines is:\n%s", lines)
	r := lines.GetTriplesForPredicate(NewStringNode("solution"))
	assert.Len(t, r.TripleSet, 1)

	// Solutions contain too many triples.  file, and line variables should only match anonymous nodes
}
