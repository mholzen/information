package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_identity(t *testing.T) {
	var x1, x2 Identity

	assert.NotEqual(t, x1, x2)

}
