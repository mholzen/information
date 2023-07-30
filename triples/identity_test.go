package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Sort(t *testing.T) {
	ten, err := NewNode(10)
	assert.Nil(t, err)
	two, err := NewNode(2)
	assert.Nil(t, err)
	assert.True(t, two.LessThan(ten))
}
