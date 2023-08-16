package transforms

import (
	"log"
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
)

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
	err := NewQuery(query, res)(tpls)
	assert.Nil(t, err)

	res, err = res.Map(NewReferences())
	assert.Nil(t, err)

	assert.Len(t, res.TripleSet, 1)

	assert.True(t, res.Contains(Triple{
		Subject:   NewStringNode("a"),
		Predicate: NewStringNode("b"),
		Object:    NewIndexNode(1),
	}))
}

func Test_NewQuery_with_multiple_conditions(t *testing.T) {
	src, err := NewJsonTriples(`[{"a":1,"b":2},{"b":2},{"c":3}]`)
	assert.Nil(t, err)

	query := NewTriples()
	x := NewVariableNode()
	query.AddTriple(x, NewStringNode("a"), NewNumberNode(1))
	query.AddTriple(x, NewStringNode("b"), NewNumberNode(2))

	res := NewTriples()
	err = src.Transform(NewQuery(query, res))
	assert.Nil(t, err)
	log.Printf("res is:\n%v", res.String())
	res, err = res.Map(NewReferences())
	assert.Nil(t, err)

	assert.Len(t, res.TripleSet, 2)
}

// func Test_NewQuery_with_unary_functions(t *testing.T) {
// 	src, err := NewJsonTriples(`[{"a":1},{"b":2},{"c":3}]`)
// 	assert.Nil(t, err)

// 	query := NewTriples()
// 	query.AddTriple(NewVariableNode(), NewStringNodeMatch("."), NewNodeMatchAny())

// 	res := NewTriples()
// 	err = src.Transform(NewQuery(query, res))

// 	assert.Nil(t, err)
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
// 	assert.Nil(t, err)

// 	res := NewTriples()
// 	err = data.Transform(NewQueryMatch(res, query))
// 	assert.Nil(t, err)
// 	assert.Len(t, res.TripleSet, 3)

// }
