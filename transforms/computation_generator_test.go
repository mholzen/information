package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/require"
)

func Test_NewObjectCountFunction(t *testing.T) {
	triples, err := NewTriplesFromStrings(
		"a b c",
		"d e c",
		"x y z",
	)
	require.Nil(t, err)

	countFunction := NewObjectCountFunction(triples)

	res, err := countFunction(tr.Str("c"))
	require.Nil(t, err)
	require.Equal(t, tr.NewIndexNode(2), res)

	res, err = countFunction(tr.Str("z"))
	require.Nil(t, err)
	require.Equal(t, tr.NewIndexNode(1), res)

	res, err = countFunction(tr.Str("x"))
	require.Nil(t, err)
	require.Equal(t, tr.NewIndexNode(0), res)

	res, err = countFunction(tr.Str("X"))
	require.Nil(t, err)
	require.Equal(t, tr.NewIndexNode(0), res)
}
