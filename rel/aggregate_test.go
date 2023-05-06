package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_aggregate_list_add(t *testing.T) {
	s1 := StringShape("foo")
	s2 := StringShape("bar")

	list := make(AggregateList[StringShape], 0)
	list.Add(s1)
	list.Add(s2)

	assert.True(t, list.Contains(s1))
	assert.True(t, list.Contains(s2))
}

func Test_aggregate_list_remove(t *testing.T) {
	s1 := StringShape("john")

	list := make(AggregateList[StringShape], 0)
	list.Add(s1)
	list.Remove(s1)
	assert.False(t, list.Contains(s1))
}

func Test_aggregate_set_add(t *testing.T) {
	set := make(AggregateSet[StringShape], 0)
	assert.False(t, set.Contains(StringShape("bing")))

	s1 := StringShape("foo")
	set.Add(s1)
	assert.True(t, set.Contains(s1))

	s2 := StringShape("bar")
	set.Add(s2)
	assert.True(t, set.Contains(s2))
}

func Test_aggregate_set_remove(t *testing.T) {
	set := make(AggregateSet[StringShape], 0)

	s1 := StringShape("john")
	set.Add(s1)
	set.Remove(s1)
	assert.False(t, set.Contains(s1))
}

func Test_aggregate_collection_pair(t *testing.T) {
	names := make(AggregateSet[StringShape], 0)

	names.Add(StringShape("john"))
	names.Add(StringShape("jane"))

	properties := make(AggregateSet[StringShape], 0)

	properties.Add(StringShape("alive"))
	properties.Add(StringShape("dead"))

	// statements := NewPairRuleSet(names, properties)

	// statement := NewPair(StringShape("john"), StringShape("alive"))
	// assert.True(statements.Contains(statement))

	// invalidStatement := NewPair(StringShape("john"), StringShape("jane"))
	// assert.False(statements.Contains(statement))
}
