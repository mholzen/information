package handlers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_ToBytes(t *testing.T) {

	fileInfo, _, err := StatWithRemainder("data/dir1/abc")
	require.Nil(t, err)

	in := Payload{
		Content: "application/json+fileinfo",
		Data:    fileInfo,
	}
	out, err := ToTextPayload(in)
	require.Nil(t, err)
	assert.Len(t, out.Data, 2)
}
