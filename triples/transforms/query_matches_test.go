package transforms

import (
	"testing"

	tr "github.com/mholzen/information/triples"
	"github.com/stretchr/testify/assert"
)

func Test_SolutionMap_TestVariables(t *testing.T) {

	v := Var()
	queryFirst := tr.NewTripleFromNodes(v, tr.Str("first"), tr.Str("Marc"))
	queryAge := tr.NewTripleFromNodes(v, tr.Str("age"), tr.Anon())

	anon := tr.Anon()
	solution1 := SolutionMap{
		queryFirst: tr.NewTripleFromNodes(anon, tr.Str("first"), tr.Str("Marc")),
		queryAge:   tr.NewTripleFromNodes(anon, tr.Str("age"), tr.NewIndexNode(18)),
	}
	assert.True(t, solution1.TestVariables())

	anon2 := tr.Anon()
	solution2 := SolutionMap{
		queryFirst: tr.NewTripleFromNodes(anon, tr.Str("first"), tr.Str("Marc")),
		queryAge:   tr.NewTripleFromNodes(anon2, tr.Str("age"), tr.Str("")),
	}
	assert.False(t, solution2.TestVariables())
}
