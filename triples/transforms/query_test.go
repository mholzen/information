package transforms

import (
	"log"
	"strings"
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_VariableNode(t *testing.T) {
	x := NewVariableNode()
	var y Node = x
	_, ok := y.(VariableNode)
	assert.True(t, ok)
}

func Test_Traverse(t *testing.T) {
	nodes := NodeList{NewStringNode("a"), NewStringNode("b"), NewStringNode("c")}
	vars := VariableList{NewVariableNode(), NewVariableNode()}
	r := vars.Traverse(nodes)
	assert.Len(t, r, 9)
}

func Test_NewQuery(t *testing.T) {
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, "b", 1)

	res := NewTriples()
	err := NewQueryTransformer(query, res, GetDefinitions())(tpls)
	require.Nil(t, err)

	require.Len(t, res.GetTriplesForPredicate(NewStringNode("solution")).TripleSet, 1)

	solutionSubject := res.GetTriplesForPredicate(NewStringNode("solution")).GetSubjectList()[0]

	solution := res.GetTriplesForSubject(solutionSubject)
	// require.Len(t, solution.TripleSet, 1)

	assert.Greater(t, len(solution.TripleSet), 1, solution.String())

	statements := solution.GetTriplesForPredicate(NewStringNode("contains"))
	require.Len(t, statements.TripleSet, 1)

	objects := statements.GetObjects()
	require.Len(t, objects, 1)

	statements = res.GetTriplesForSubjects(objects)
	require.Len(t, statements.TripleSet, 3)

	references, err := NewReferences()(statements)
	require.Nil(t, err)

	assert.Len(t, references.TripleSet, 1)

	assert.True(t, references.Contains(Triple{
		Subject:   NewStringNode("a"),
		Predicate: NewStringNode("b"),
		Object:    NewIndexNode(1),
	}))
}

func Test_NewQuery_with_multiple_conditions(t *testing.T) {
	src, err := NewJsonTriples(`[{"a":1,"b":2},{"b":2},{"c":3}]`)
	require.Nil(t, err)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, NewStringNode("a"), NewFloatNode(1))
	query.AddTriple(x, NewStringNode("b"), NewFloatNode(2))

	res := NewTriples()
	err = src.Transform(NewQueryTransformer(query, res, GetDefinitions()))
	require.Nil(t, err)
	log.Printf("res is:\n%v", res.String())
	res, err = res.Map(NewReferences())
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
	tpls := NewTriples()
	tpls.AddTriple("a", "b", 1)
	tpls.AddTriple("b", "c", 2)
	tpls.AddTriple("e", "f", 3)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, NewNodeMatchAny(), 2)

	res := NewTriples()
	err := NewQueryTransformer(query, res, GetDefinitions())(tpls)
	require.Nil(t, err)
	assert.Greater(t, len(res.TripleSet), 1) // ContainsTriples modifies the triples

	res, err = res.Map(NewReferences())
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
	// Load a CSV file
	tm := NewCsvParser(strings.NewReader("a,b,c\nd,e,f"))

	logrus.SetLevel(logrus.DebugLevel)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)

	log.Printf("src is:\n%s", src)

	query := NewTriples()
	file := NewVariableNode()
	query.AddTriple(file, "type", "triples.AnonymousNode")
	query.AddTriple(file, "source", "CsvParser")

	// query for lines
	lines := NewTriples()
	err = NewQueryTransformer(query, lines, GetDefinitions())(src)
	require.Nil(t, err)

	// Count solutions
	log.Printf("lines is:\n%s", lines)
	r := lines.GetTriplesForPredicate(NewStringNode("solution"))
	assert.Len(t, r.TripleSet, 1)

	// Solutions contain too many triples.  file, and line variables should only match anonymous nodes
}
