package set

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_parse(t *testing.T) {
	statements, err := parse("test.txt")
	assert.Nil(t, err)
	log.Printf("%+v", statements)
}

func Test_parse2(t *testing.T) {
	statements, err := parse("a b c")
	assert.Nil(t, err)
	assert.Contains(t, statements, NewStatement("a", "b", "c"))
}
