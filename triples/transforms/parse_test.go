package transforms

import (
	"log"
	"strings"
	"testing"

	. "github.com/mholzen/information/triples"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_remove_comment(t *testing.T) {
	assert.Empty(t, removeComment("// comment"))
	assert.Equal(t, "abc ", removeComment("abc // comment"))
	assert.Equal(t, `"http://a.c" `, removeComment(`"http://a.c" // comment`))
}

func Test_parse(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}`)

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.NotNil(t, tm.Result)

	assert.Len(t, src.TripleSet, 2)
}

func Test_triples_map2(t *testing.T) {
	tm := NewJsonParser(`{"first":"marc","last":"von Holzen"}]`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.Len(t, src.TripleSet, 2)
}

func Test_triples_array_object(t *testing.T) {
	tm := NewJsonParser(`{"names":["marc","Marc", "Marco"]}`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.Len(t, src.TripleSet, 4)
}

func Test_slice_as_object(t *testing.T) {
	tm := NewJsonParser(`["root", "contains", ["marc", "is", "alive"]]`)
	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)
	assert.Len(t, src.TripleSet, 6)
}

func Test_NewFileJsonParser(t *testing.T) {
	tm := NewFileJsonParser("../../data/object.jsonc")

	dest := NewTriples()
	err := dest.Transform(tm.Transformer)
	require.Nil(t, err)

	assert.Greater(t, len(dest.TripleSet), 10)
}

func Test_csv_parse(t *testing.T) {
	tm := NewCsvParser(strings.NewReader("a,b,c\nd,e,f"))

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)

	assert.Len(t, src.TripleSet, 9)
}

func Test_NewLinesParser(t *testing.T) {
	tm := NewLinesParser(strings.NewReader("line 1\nline 2\nline 3\n"))

	src := NewTriples()
	err := src.Transform(tm.Transformer)
	require.Nil(t, err)

	assert.Len(t, src.TripleSet, 4)
}

func Test_NewNodeFromString(t *testing.T) {
	node := NewNodeFromString("marc")
	require.Equal(t, node, NewStringNode("marc"))

	node = NewNodeFromString("?")
	require.IsType(t, NewVariableNode(), node)

	node = NewNodeFromString("_")
	require.IsType(t, NewAnonymousNode(), node)

	node = NewNodeFromString("42")
	require.Equal(t, node, NewIndexNode(42))

	node = NewNodeFromString("42.0")
	require.Equal(t, node, NewFloatNode(42.0))
}

func Test_NewTripleFromString(t *testing.T) {
	triple, err := NewTripleFromString(
		"_ first Marc",
	)
	require.Nil(t, err)
	require.IsType(t, NewAnonymousNode(), triple.Subject)
	require.Equal(t, triple.Predicate, NewStringNode("first"))
	require.Equal(t, triple.Object, NewStringNode("Marc"))
}

func Test_NewTriplesFromtStrings(t *testing.T) {
	tpls, err := NewTriplesFromStrings(
		"marc first Marc",
		"? first ?",
		"? age ?",
	)
	require.Nil(t, err)

	log.Printf("%+v", tpls)
	assert.Len(t, tpls.TripleSet, 3)
}

func Test_NewNamedTriples(t *testing.T) {
	query, err := NewNamedTriples(
		"_n foo ?x",
		"?x bar _n",
	)
	require.Nil(t, err)
	require.Len(t, query.TripleSet, 2)
	triples := query.GetTripleList()
	require.Equal(t, triples[0].Subject, triples[1].Object)
	require.Equal(t, triples[1].Subject, triples[0].Object)
}

func Test_NewNamedTriples_Function(t *testing.T) {
	t1 := NewTripleFromNodes(TypeFunctionNode, Str("name"), Str("type"))
	names := NewNamedNodeMapFromTriples(
		NewTriples().Add(t1),
	)
	query, err := names.NewTriples(
		"?x type() string",
	)
	require.Nil(t, err)
	require.Len(t, query.TripleSet, 1)

	triples := query.GetTripleList()
	f, ok := triples[0].Predicate.(UnaryFunctionNode)
	require.True(t, ok)
	val, err := f(Str("foo"))
	require.Nil(t, err)

	require.Equal(t, "triples.StringNode", val.String())
}
