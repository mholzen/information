package rel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_compute(t *testing.T) {

	c := compute()
	assert.Greater(t, c, 0)
}
