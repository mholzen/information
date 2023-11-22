package transforms

import (
	"log"
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_RowQuery(t *testing.T) {
	rowQuery := RowQuery()
	source := array()

	res, err := source.Map(NewQueryMapper(rowQuery))
	require.Nil(t, err)

	res = References(res)
	require.Len(t, res.TripleSet, 2)
}

func Test_RowTriples(t *testing.T) {
	source := array()

	res, err := RowTriples(source)
	require.Nil(t, err)

	require.Len(t, res.TripleSet, 2)
}

func array() *tr.Triples {
	source := tr.NewTriples()
	source.AddTriple(tr.NewAnonymousNode(), 1, tr.NewAnonymousNode())
	source.AddTriple(tr.NewAnonymousNode(), 2, tr.NewAnonymousNode())
	source.AddTriple(tr.NewAnonymousNode(), "c", tr.NewAnonymousNode())
	return source
}

func matrix() *tr.Triples {
	res := tr.NewTriples()
	root := tr.NewAnonymousNode()
	r1 := tr.NewAnonymousNode()
	r2 := tr.NewAnonymousNode()

	res.AddTriple(root, 1, r1)
	res.AddTriple(root, 2, r2)
	res.AddTriple(r1, 1, "11")
	res.AddTriple(r1, 2, "12")
	res.AddTriple(r2, 1, "21")
	res.AddTriple(r2, 2, "22")

	log.Printf("matrix: %s", res)
	return res
}

func Test_MatrixQuery(t *testing.T) {
	query := MatrixQuery()
	res, err := matrix().Map(NewQueryMapper(query))
	require.Nil(t, err)

	res = References(res)
	require.Len(t, res.TripleSet, 2)
}
