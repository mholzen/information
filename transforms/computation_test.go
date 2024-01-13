package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_NewComputation(t *testing.T) {
	v := Var()
	computation := NewComputation(v, tr.LengthFunctionNode, tr.NewIndexNode(4))

	r, err := computation.Test(tr.Str("Joe"))
	require.Nil(t, err)
	assert.False(t, r)

	r, err = computation.Test(tr.Str("Marc"))
	require.Nil(t, err)
	assert.True(t, r)
}

func Test_Computations(t *testing.T) {
	v := Var()
	lengthEquals4 := NewComputation(v, tr.LengthFunctionNode, tr.NewIndexNode(4))
	containsAr := NewComputation(v, tr.NewStringNodeMatch("ar"), tr.NewNumberNode(1))
	computations := NewComputations(lengthEquals4, containsAr)

	assert.False(t, computations.Test(VariableMap{v: tr.Str("Joe")}))
	assert.True(t, computations.Test(VariableMap{v: tr.Str("Marc")}))
}
