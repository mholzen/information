package triples

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func decodeJson(input string) (interface{}, error) {
	buffer := bytes.NewBuffer([]byte(input))
	decoder := json.NewDecoder(strings.NewReader(buffer.String()))
	var data interface{}
	err := decoder.Decode(&data)
	return data, err
}

func Test_parse2(t *testing.T) {
	data, err := decodeJson(`{"first":"marc","last":"von Holzen"}`)
	assert.Nil(t, err)

	triples := TriplesModifier{Triples: NewTriples()}
	node, err := triples.Parse(data)
	assert.Nil(t, err)
	assert.Len(t, triples.TripleSet, 2)
	assert.NotNil(t, node)
}
