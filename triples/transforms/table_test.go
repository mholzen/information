package transforms

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewTableDefinitionSmall(t *testing.T) {
	src, err := NewJsonTriples(`[
		{"first":"marc"}
		]`)
	require.Nil(t, err)
	log.Printf("src: %s", src)

	res, err := src.Map(TableMapper)
	require.Nil(t, err)
	log.Printf("res: %s", res)
	assert.Equal(t, 4, len(res.TripleSet))
}
