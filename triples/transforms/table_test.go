package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTableDefinition(t *testing.T) {
	def := triples.NewTriples()
	def.AddTriple("x", 0, "first")
	def.AddTriple("x", 1, "last")
	def.AddTriple("x", 2, "age")
	def.AddTriple("x", 3, "first")
	def.AddTriple("x", 3, "last")
	col := NewTableDefinition(def)
	assert.Equal(t, 4, len(col.Columns))
	assert.Equal(t, 2, len(col.Columns[3]))

	src, err := NewJsonTriples(`[
		{"first":"marc","last":"von Holzen", "age": 52},
		{"first":"John","last":"Doe", "age": 22},
		{"first":"Jane","last":"Wilkenson", "age": 28}
		]`)
	require.Nil(t, err)

	table := NewTableGenerator(def)
	err = src.Transform(table.Transformer)
	require.Nil(t, err)
	assert.Equal(t, 4, len(table.Rows))
	assert.Equal(t, 4, len(table.Rows[0]))
}
