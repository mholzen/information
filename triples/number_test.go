package triples

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CompareNumbers(t *testing.T) {
	nn1 := NewNumberNode(1)
	nn2 := NewNumberNode(1)
	assert.True(t, nn1 == nn2)
}
