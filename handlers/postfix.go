package handlers

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"
	"github.com/mholzen/information/triples/transforms/html"
)

func StatWithRemainder(filePath string) (FileInfo, []string, error) {
	filePath = path.Clean(filePath)
	components := strings.Split(filePath, "/")

	currentPath := os.Getenv("ROOT") // TODO: put this in a factory
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

func GetHandlerMap() (map[string]Transform, error) {
	// load mappers
	rowMapper, err := transforms.RowMapper()
	if err != nil {
		return nil, err
	}

	res := map[string]Transform{
		"content":             ToContent,
		"graph":               ToGraphPayload,
		"html":                ToHtml,
		"list":                ToListPayload,
		"style":               ToStylePayload,
		"mime":                ToMimeType,
		"nodelink":            ToNodeLinkPayload,
		"tableDefinition":     ToTableDefinitionPayload,
		"predicates":          ToTableDefinitionPayload,
		"text":                ToTextPayload,
		"triples":             ToTriplesPayload,
		"transform,rows":      NewMapperPayload(rowMapper),
		"transform,table":     NewMapperPayload(transforms.TableMapper),
		"transform,htmlTable": NewMapperPayload(html.HtmlTableMapper),
		"transform,id,lines":  NewMapperPayload(transforms.IdLinesMapper),
		"data":                ToDataPayload,
	}
	return res, nil
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

	handlerMap, err := GetHandlerMap()
	if err != nil {
		return err
	}

	payload := Payload{
		Content: "application/json+fileinfo",
		Data:    fileInfo,
	}
	for _, name := range remainder {
		if handler, ok := handlerMap[name]; !ok {
			return echo.NewHTTPError(404, fmt.Sprintf("'%s' handler not Found", name))
		} else {
			payload, err = handler(payload)
			if err != nil {
				return err
			}
		}
	}

	if s, ok := payload.Data.(triples.Node); ok {
		return c.String(200, s.String())
	} else if strings.HasPrefix(payload.Content, "text/html") {
		return c.HTML(200, payload.Data.(string))
	} else if strings.HasPrefix(payload.Content, "text/") {
		return c.String(200, payload.Data.(string))
	} else if strings.HasPrefix(payload.Content, "application/json") {
		return c.JSON(200, payload.Data)
	} else if s, ok := payload.Data.(string); ok {
		return c.String(200, s)
	} else {
		return c.Blob(200, payload.Content, payload.Data.([]byte))
	}
}
