package html

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_html(t *testing.T) {
	tm := transforms.NewFileJsonParser("../../data/examples.jsonc")
	// all, err := Parse("../data/examples.jsonc")
	res := triples.NewTriples()
	err := res.Transform(tm.Transformer)
	require.Nil(t, err)

	html := NewHtmlTransformer(*res, res.GetTripleList(), 4)
	assert.Greater(t, len(html.String()), 1000)
}
