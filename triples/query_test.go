package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_variable_node(t *testing.T) {
	data, err := DecodeJson(`[{"a":1},{"b":2},{"c":3}]`) // TODO: this 1 is created as a Float node
	assert.Nil(t, err)

	src := NewTriples()
	err = src.Transform(NewParser(data))
	assert.Nil(t, err)

	query := NewTriples()
	query.AddTriple(NewVariableNode(), "a", 1.0)

	res := NewTriples()
	err = src.Transform(NewQueryMatch(res, query))
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 1)
}

func Test_match(t *testing.T) {
	tpls := NewTriples()
	x := NewVariableNode()
	value, _ := tpls.AddTriple("a", "b", 1)
	query, _ := tpls.AddTriple(x, "b", 1)
	assert.True(t, Match(value, query))
	assert.False(t, x.Excludes.Contains(NewStringNode("a")))

	query2, _ := tpls.AddTriple(x, "b", 2)
	assert.False(t, Match(value, query2))
	assert.True(t, x.Excludes.Contains(NewStringNode("a")))

}
