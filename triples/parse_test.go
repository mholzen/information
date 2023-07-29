package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_remove_comment(t *testing.T) {
	assert.Empty(t, RemoveComment("// comment"))
	assert.Equal(t, "abc ", RemoveComment("abc // comment"))
	assert.Equal(t, `"http://a.c" `, RemoveComment(`"http://a.c" // comment`))
}

func Test_parse_old(t *testing.T) {
	all, err := Parse("../data/examples.jsonc")

	if !assert.Nil(t, err) {
		return
	}

	subject := NewStringNode("marc")
	marc := all.AddReachableTriples(subject, nil)

	objectTriples := marc.GetTriplesForObject(subject, nil)
	assert.Greater(t, len(objectTriples.TripleSet), 10)
}
