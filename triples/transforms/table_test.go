package transforms

import (
	"log"
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
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

	data := triples.NewJsonParser(`[
		{"first":"marc","last":"von Holzen", "age": 52},
		{"first":"John","last":"Doe", "age": 22},
		{"first":"Jane","last":"Wilkenson", "age": 28}
		]`)

	src := triples.NewTriples()
	err := src.Transform(data.Transformer)
	assert.Nil(t, err)

	table := NewTableGenerator(def)
	err = src.Transform(table.Transformer)
	assert.Nil(t, err)
	log.Println(table.Rows[0])
	log.Println(table.Rows[1])
	log.Println(table.Rows[2])
	log.Println(table.Rows[3])
	assert.Equal(t, 12, len(table.Rows)) // TODO: should be 3 -- is 12 because all triples are evaluated
	assert.Equal(t, 4, len(table.Rows[0]))
}
