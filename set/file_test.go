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
