package handlers

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
)

func Test_convert(t *testing.T) {
	v := MimeType("text/plain")
	var node triples.Node = v
	assert.Equal(t, node.String(), "text/plain")
}
