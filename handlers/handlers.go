package handlers

import (
	"bufio"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/mholzen/information/triples"
	"github.com/mholzen/information/triples/transforms"
	"github.com/mholzen/information/triples/transforms/node_link"
)

func Filepath(url string) string {
	return filepath.Join(os.Getenv("ROOT"), url)
}

func Extension(path string) string {
	return filepath.Ext(path)
}

func Content(path string) (io.Reader, error) {
	return transforms.ReadAndStripComments(path)
}

func Triples(url string) (*triples.Triples, error) {
	data, err := Content(Filepath(url))
	if err != nil {
		return nil, err
	}
	mimeType := mime.TypeByExtension(Extension(url))
	if mimeType == "" {
		reader := bufio.NewReader(data)
		peeked, err := reader.Peek(512)
		if err != io.EOF && err != nil {
			return nil, err
		}
		mimeType = http.DetectContentType(peeked)

		// reading on the bufio.Reader will advance the reader nevertheless, so re-assign data to the new reader
		data = reader
	}

	tm, err := transforms.NewParserFromContentType(mimeType, data)
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

func FilesHandler(c echo.Context) error {
	Content(Filepath(c.Param("file")))
	return nil
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

func StatsHandler(c echo.Context) error {
	path := Filepath(c.Param("file"))
	stat, err := os.Stat(path)
	if os.IsNotExist(err) {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	if err != nil {
		return err
	}
	if stat.IsDir() {
		entries, err := os.ReadDir(path)
		if err != nil {
			return err
		}
		list := []FileInfo{}

		for _, entry := range entries {

			stat, err := os.Stat(filepath.Join(path, entry.Name()))

			f := FileInfo{
				Name:    entry.Name(),
				Size:    stat.Size(),
				Mode:    stat.Mode(),
				ModTime: stat.ModTime(),
				IsDir:   entry.IsDir(),
				Error:   err,
			}
			list = append(list, f)
		}

		return c.JSON(http.StatusOK, list)
	} else {
		f := FileInfo{
			Name:    stat.Name(),
			Size:    stat.Size(),
			Mode:    stat.Mode(),
			ModTime: stat.ModTime(),
			IsDir:   stat.IsDir(),
		}

		return c.JSON(http.StatusOK, f)
	}
}

type FileInfo struct {
	Name    string
	Size    int64
	Mode    os.FileMode
	ModTime time.Time
	IsDir   bool
	Error   error
}

func HtmlHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return err
	}
	tripleList := src.GetTripleList()
	tripleList.Sort()

	html := transforms.NewHtmlTransformer(*src, tripleList, 0)

	data := map[string]interface{}{
		"tripleList": tripleList,
		"html":       html.String(),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func ObjectsHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return err
	}

	dest := triples.NewAnonymousNode()
	err = src.Transform(transforms.NewTraverse(dest, transforms.AlwaysTripleMatch, dest, src))
	if err != nil {
		return err
	}

	dest2 := triples.NewAnonymousNode()
	objectMapper := transforms.NewTripleObjectTransformer(dest2, src)
	err = src.Transform(transforms.NewMap(dest, objectMapper, src))
	if err != nil {
		return err
	}

	res := triples.NewTriples()
	err = src.Transform(transforms.NewFlatMap(dest2, transforms.GetStringObjectMapper, res))
	if err != nil {
		return err
	}

	answer := res.GetTripleList().GetObjectStrings()
	sort.Strings(answer)
	return c.String(http.StatusOK, strings.Join(answer, "\n"))
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

func TableHandler(c echo.Context) error {
	src, err := Triples(c.Param("file"))
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	def := transforms.NewDefaultTableDefinition(src)
	tr := transforms.NewTableGenerator(def)
	err = src.Transform(tr.Transformer)
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.HTML(http.StatusOK, tr.Html())
}