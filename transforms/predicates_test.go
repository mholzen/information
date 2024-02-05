package transforms

import (
	"testing"

	"github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_PredicatesSortedLexical(t *testing.T) {
	tpls, err := NewTriplesFromStrings(
		"a b c",
		"d e f",
		"g h i",
	)
	require.Nil(t, err)

	preds := PredicatesSortedLexical(tpls)
	assert.Equal(t, "b", preds.GetTriplesForPredicate(triples.NewIndexNode(0)).GetObjectList()[0].String())
	assert.Equal(t, "e", preds.GetTriplesForPredicate(triples.NewIndexNode(1)).GetObjectList()[0].String())
	assert.Equal(t, "h", preds.GetTriplesForPredicate(triples.NewIndexNode(2)).GetObjectList()[0].String())
}
