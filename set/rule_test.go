package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_aggregate_collection_pair(t *testing.T) {
	rule := getRuleList()

	statement := NewList([]StringShape{"john", "alive"}...)
	assert.True(t, rule.Contains(statement))

	invalidStatement := NewList([]StringShape{"john", "jane"}...)
	assert.False(t, rule.Contains(invalidStatement))

	invalidStatement2 := NewList([]StringShape{"john", "alive", "well"}...)
	assert.False(t, rule.Contains(invalidStatement2))
}

func getRuleList() ListRuleList[StringShape] {
	subjects := NewMapSet[StringShape]()
	subjects.Add("john")
	subjects.Add("jane")

	properties := NewMapSet[StringShape]()
	properties.Add("alive")
	properties.Add("dead")

	statements := NewListRuleList[StringShape](subjects, properties)
	return statements
}

func Test_validate_statement_lists(t *testing.T) {

	statement1 := NewList([]StringShape{"john", "alive"}...)
	statement2 := NewList([]StringShape{"john", "xxx"}...)

	statements := NewList(statement1, statement2)

	rule := getRuleList()

	valid := statements.Intersect(rule)
	assert.Equal(t, len(valid), 1)

	invalid := statements.Difference(rule)
	assert.Equal(t, len(invalid), 1)
}
