package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_html(t *testing.T) {
	tm := NewFileJsonParser("../data/examples.jsonc")
	// all, err := Parse("../data/examples.jsonc")
	res := NewTriples()
	err := res.Transform(tm.Transformer)
	assert.Nil(t, err)

	html := NewHtmlTransformer(*res, res.GetTripleList(), 4)
	assert.Greater(t, len(html.String()), 1000)
}
