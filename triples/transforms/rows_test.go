package transforms

import (
	"log"
	"os"
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_RowQuery(t *testing.T) {
	query := RowQuery2()
	source, err := NewTriplesFromStrings(
		"_ 1 _",
		"_ 2 _",
		"_ 3 _",
		"_ 4 1",
		"a 5 _",
	)
	require.Nil(t, err)

	solutions, err := query.Apply(source)
	require.Nil(t, err)

	log.Printf("%s", solutions.GetAllTriples())
	require.Len(t, solutions, 3)
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
	return res
}

func TestMain(m *testing.M) {
	// Initialize logging for all tests
	// (Set up log file, log format, etc.)
	InitLog()
	// Run the tests
	os.Exit(m.Run())
}

func Test_MatrixQuery2(t *testing.T) {
	t.Skip()
	// query, err := NewQueryFromTriples(MatrixQuery())
	// require.Nil(t, err)

	// require.Len(t, query.QueryTriples, 2)

	// MatrixQuery triples don't have an order, so I can't tell which triple gives me the rows
	// must separate graph selection from graph condition
	// could that be two queries?

	// sol, err := query.GetSolutions(matrix())
	// require.Nil(t, err)

	// require.Len(t, sol.SelectTriples().TripleSet, 2)
}

func Test_MatrixQuery3(t *testing.T) {
	query := MatrixQuery()
	selects := ReferenceTriplesConnected(query)
	log.Printf("query:\n%s", query)
	require.Len(t, selects.TripleSet, 4)
}
