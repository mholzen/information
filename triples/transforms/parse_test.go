package transforms

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_remove_comment(t *testing.T) {
	assert.Empty(t, RemoveComment("// comment"))
	assert.Equal(t, "abc ", RemoveComment("abc // comment"))
	assert.Equal(t, `"http://a.c" `, RemoveComment(`"http://a.c" // comment`))
}
