package handlers

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/information/transforms"
	"github.com/mholzen/information/transforms/html"
	"github.com/mholzen/information/transforms/node_link"
	"github.com/mholzen/information/triples"
)

func Filepath(url string) string {
	return filepath.Join(os.Getenv("ROOT"), url)
}

func Extension(path string) string {
	return filepath.Ext(path)
}

func Triples(url string) (*triples.Triples, error) {
	reader, err := transforms.ReadAndStripComments(url)
	if err != nil {
		return nil, err
	}
	mimeType := mime.TypeByExtension(Extension(url))
	if mimeType == "" {
		newReader := bufio.NewReader(reader)
		peeked, err := newReader.Peek(512)
		if err != io.EOF && err != nil {
			return nil, err
		}
		mimeType = http.DetectContentType(peeked)

		// reading on the bufio.Reader will advance the reader nevertheless, so re-assign data to the new reader
		reader = newReader
	}

	tm, err := transforms.NewParserFromContentType(mimeType, reader)
	if err != nil {
		return nil, err
	}
	res := triples.NewTriples()
	err = res.Transform(tm.Transformer)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func TriplesHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return err
	}
	tripleList := src.GetTripleList()
	tripleList.Sort()

	return c.JSON(http.StatusOK, tripleList)
}

func HtmlHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return err
	}
	tripleList := src.GetTripleList()
	tripleList.Sort()

	html := html.NewHtmlTransformer(*src, tripleList, 0)

	data := map[string]interface{}{
		"tripleList": tripleList,
		"html":       html.String(),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func NodeLinkHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return err
	}
	tr := node_link.NewNodeLinkTransformer()
	err = src.Transform(tr.Transformer)
	if err != nil {
		return err
	}
	return c.JSON(http.StatusOK, tr.Result)
}

func GraphHandler(c echo.Context) error {
	data := map[string]interface{}{
		"url": "/nodelink/" + c.Param("file"),
	}
	return c.Render(http.StatusOK, "graph.html", data)
}

func DataHandler(c echo.Context) error {
	data := map[string]interface{}{
		"url": "/nodelink/" + c.Param("file"),
	}
	return c.Render(http.StatusOK, "graph.html", data)
}
