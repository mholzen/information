package triples

import (
	"bytes"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	all, err := Parse("examples.json")
	assert.Nil(t, err)

	subject := NewStringNode("marc")
	marc := all.AddReachableTriples(subject, nil)

	subjectTriples := marc.GetTriplesForSubject(subject, nil)
	log.Printf("=== subject triples ===")
	res := fmt.Sprintf("%s\n", subject)
	res += marc.String(subjectTriples, "   ", 4)
	log.Printf("\n%s", res)

	objectTriples := marc.GetTriplesForObject(subject, nil)
	log.Printf("=== object triples for marc")
	for triple := range objectTriples.TripleSet {
		log.Printf("%s", triple)
	}
}

func Test_triple_in_subject_predicate_object(t *testing.T) {
	buffer := bytes.NewBuffer([]byte(`{"s":"marc","p":"is","o":"alive"}`))
	parseString(*buffer)
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
