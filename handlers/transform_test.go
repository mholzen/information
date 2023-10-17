package handlers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ToBytes(t *testing.T) {

	os.Setenv("ROOT", "../")
	fileInfo, _, err := StatWithRemainder("data/dir1/content")
	require.Nil(t, err)

	in := Payload{
		Content: "application/json+fileinfo",
		Data:    fileInfo,
	}
	out, err := ToText(in)
	require.Nil(t, err)
	assert.Len(t, out.Data, 1)
}
