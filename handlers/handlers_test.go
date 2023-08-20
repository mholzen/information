package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func Test_File(t *testing.T) {
	e := echo.New()

	req, err := http.NewRequest("POST", "/echo", nil)
	assert.Nil(t, err)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/files/:file")
	c.SetParamNames("file")
	c.SetParamValues("foo")

	if assert.NoError(t, FilesHandler(c)) {
		assert.Equal(t, rec.Code, http.StatusOK)
		assert.Equal(t, strings.Trim(rec.Body.String(), "\n"), "foo")
	}
}
