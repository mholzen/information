package handlers

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_File(t *testing.T) {
	e := echo.New()
	os.Setenv("ROOT", "../")

	req, err := http.NewRequest(http.MethodGet, "/stats/data/object.jsonc", nil)
	require.Nil(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/stats/:file")
	c.SetParamNames("file")
	c.SetParamValues("data/object.jsonc")

	if assert.NoError(t, StatsHandler(c)) {
		assert.Equal(t, rec.Code, http.StatusOK)
	}
}
