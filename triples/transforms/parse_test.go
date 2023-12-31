package transforms

import (
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
