package transforms

import (
	"testing"

	. "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_replace(t *testing.T) {
	data := NewTriples()
	data.AddTripleString("a", "b", "c")
	data.AddTripleString("d", "e", "f")
	x := NewVariableNode()
	data.AddTripleFromAny(x, "h", "i")
	replace := NewReplaceMapper(NewStringNode("a"), NewStringNode("A"))
	res, err := data.Map(replace)
	require.Nil(t, err)

	assert.True(t, res.Contains(Triple{
		Subject:   NewStringNode("A"),
		Predicate: NewStringNode("b"),
		Object:    NewStringNode("c"),
	}))

	replace = NewReplaceMapper(x, NewStringNode("X"))
	res, err = data.Map(replace)
	require.Nil(t, err)

	assert.True(t, res.Contains(Triple{
		Subject:   NewStringNode("X"),
		Predicate: NewStringNode("h"),
		Object:    NewStringNode("i"),
	}))

}
