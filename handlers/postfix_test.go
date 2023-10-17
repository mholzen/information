package handlers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_StatFile(t *testing.T) {
	if os.Getenv("ROOT") == "" {
		os.Setenv("ROOT", "../")
	}

	fileInfo, remainder, err := StatWithRemainder("data/index.md/html")

	require.Nil(t, err)
	assert.Equal(t, []string{"html"}, remainder)
	assert.Equal(t, "../data/index.md", fileInfo.Name)
}
