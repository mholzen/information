package transforms

import (
	"log"
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_Root(t *testing.T) {
	src := triples.NewTriples()
	a := triples.NewAnonymousNode()
	b := triples.NewAnonymousNode()

	src.AddTriple(b, "c", "c")
	src.AddTriple(a, "b", b)

	res, err := Root(src)
	require.Nil(t, err)

	log.Printf("res: %s", res)
	assert.Equal(t, a, res)
}
