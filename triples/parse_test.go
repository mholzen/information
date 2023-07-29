package triples

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_remove_comment(t *testing.T) {
	assert.Empty(t, RemoveComment("// comment"))
	assert.Equal(t, "abc ", RemoveComment("abc // comment"))
	assert.Equal(t, `"http://a.c" `, RemoveComment(`"http://a.c" // comment`))
}

func Test_parse(t *testing.T) {
	all, err := Parse("../data/examples.jsonc")

	if !assert.Nil(t, err) {
		return
	}

	subject := NewStringNode("marc")
	marc := all.AddReachableTriples(subject, nil)

	objectTriples := marc.GetTriplesForObject(subject, nil)
	assert.Greater(t, len(objectTriples.TripleSet), 10)
}

func Test_triple_in_subject_predicate_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`{"s":"marc","p":"is","o":"alive"}`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 3)
}

func Test_triples_map(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`{"first":"marc","last":"von Holzen"}]`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 2)
}

func Test_triples_array_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`{"names":["marc","Marc", "Marco"]}`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 4)
}

func Test_slice_as_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`["root", "contains", ["marc", "is", "alive"]]`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 6)
}
