package handlers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

func Split(filePath string) []string {
	filePath = path.Clean(filePath)
	return strings.Split(filePath, "/")
}

func StatWithRemainder(path string) (FileInfo, []string, error) {
	components := Split(path)
	currentPath := os.Getenv("ROOT")
	var currentFileInfo os.FileInfo
	for i, component := range components {
		nextPath := filepath.Join(currentPath, component)
		if nextFileInfo, err := os.Stat(nextPath); os.IsNotExist(err) {
			// If path doesn't exist, return the current path and the remaining path
			res := NewFileInfo(currentFileInfo)
			res.Name = currentPath
			return res, components[i:], nil
		} else if !nextFileInfo.IsDir() {
			res := NewFileInfo(nextFileInfo)
			res.Name = nextPath

			return res, components[i+1:], nil
		} else {
			currentPath = nextPath
			currentFileInfo = nextFileInfo
		}
	}
	res := NewFileInfo(currentFileInfo)
	res.Name = currentPath
	return res, nil, nil
}

var HandlerMap = map[string]Transform{
	"html":    ToHtml,
	"content": ToText,
}

func FilesPostfixHandler(c echo.Context) error {
	// Warning: at every path segment, we may want to serve the directory if it
	// exists, or serve the handler if it exists which introduces indeterminism

	// Given /files/data/html/foo.html, does that first `html` refer to a
	// directry in data or the handler named `html`?

	arg := c.Param("file")

	fileInfo, remainder, err := StatWithRemainder(arg)
	if err != nil {
		return err
	}

	payload := Payload{
		Content: "application/json+fileinfo",
		Data:    fileInfo,
	}
	for _, name := range remainder {
		if handler, ok := HandlerMap[name]; !ok {
			return echo.NewHTTPError(404, fmt.Sprintf("'%s' handler not Found", name))
		} else {
			payload, err = handler(payload)
			if err != nil {
				return err
			}
		}
	}

	if strings.HasPrefix(payload.Content, "text/html") {
		return c.HTML(200, payload.Data.(string))
	} else if strings.HasPrefix(payload.Content, "application/json") {
		return c.JSON(200, payload.Data)
	} else if s, ok := payload.Data.(string); ok {
		return c.String(200, s)
	} else {
		return c.Blob(200, payload.Content, payload.Data.([]byte))
	}
}
