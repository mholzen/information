package triples

import (
	"bytes"
	"log"
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

	log.Printf("=== triples ===")
	log.Printf("%s\n", marc.String())

	objectTriples := marc.GetTriplesForObject(subject, nil)
	log.Printf("=== object triples for marc")
	for triple := range objectTriples.TripleSet {
		log.Printf("%s", triple)
	}
}

func Test_triple_in_subject_predicate_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`{"s":"marc","p":"is","o":"alive"}`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 6)
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
	assert.Len(t, res.TripleSet, 3)
}

func Test_slice_as_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`["root", "contains", ["marc", "is", "alive"]]`))
	res, err := parseString(*buffer)
	assert.Nil(t, err)
	assert.Len(t, res.TripleSet, 6)
}

func Test_html(t *testing.T) {
	all, err := Parse("examples.json")
	assert.Nil(t, err)

	subject := NewStringNode("marc")
	marc := all.AddReachableTriples(subject, nil)

	subjectTriples := marc.GetTriplesForSubject(subject, nil)
	html := NewHtmlTransformer(*marc, subjectTriples, 4)
	log.Printf("=== subject triples ===\n%s", html.String())
}
