package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NodeLexicalSort(t *testing.T) {
	nodes, err := NewNodeList(
		"z",
		"a",
		"Z",
		"A",
		false,
		true,
		2,
		11,
		1.1,
		1,
	)
	require.Nil(t, err)

	nodes.SortLexical()

	assert.Equal(t, []string{
		"true",
		"false",
		"1",
		"1.100000",
		"2",
		"11",
		"A",
		"Z",
		"a",
		"z",
	}, nodes.Strings())
}
