package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_NewTriplesFromMap(t *testing.T) {
	triples := NewTriples()
	_, err := triples.NewTriplesFromMap(map[string]interface{}{
		"first": "marc",
		"last":  "von Holzen",
		"age":   50,
	})
	assert.Nil(t, err)
	assert.Len(t, triples.TripleSet, 3)
}
