package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_html(t *testing.T) {
	all, err := Parse("../data/examples.jsonc")
	assert.Nil(t, err)

	html := NewHtmlTransformer(*all, all.GetTripleList(), 4)
	assert.Greater(t, len(html.String()), 1000)
}