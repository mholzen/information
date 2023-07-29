package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"

	"github.com/mholzen/information/triples"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func GetTripleList(c echo.Context) (*triples.Triples, error) {
	root := os.Getenv("ROOT")
	path := filepath.Join(root, c.Param("file"))
	all, err := triples.Parse(path)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return all, nil
}

func HtmlHandler(c echo.Context) error {
	src, err := GetTripleList(c)
	if err != nil {
		return err
	}
	tripleList := src.GetTripleList()
	tripleList.Sort()

	html := triples.NewHtmlTransformer(*src, tripleList, 0)

	data := map[string]interface{}{
		"tripleList": tripleList,
		"html":       html.String(),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func ObjectsHandler(c echo.Context) error {
	src, err := GetTripleList(c)
	if err != nil {
		return err
	}

	dest := triples.NewAnonymousNode()
	err = src.Transform(triples.NewTraverse(dest, triples.AlwaysTripleMatch, dest, src))
	if err != nil {
		return err
	}

	dest2 := triples.NewAnonymousNode()
	objectMapper := triples.NewTripleObjectTransformer(dest2, src)
	err = src.Transform(triples.NewMap(dest, objectMapper, src))
	if err != nil {
		return err
	}

	res := triples.NewTriples()
	err = src.Transform(triples.NewFlatMap(dest2, triples.GetStringObjectMapper, res))
	if err != nil {
		return err
	}

	answer := res.GetTripleList().GetObjectStrings()
	sort.Strings(answer)
	return c.String(http.StatusOK, strings.Join(answer, "\n"))

}

func main() {
	e := echo.New()
	e.GET("/html/:file", HtmlHandler)
	e.GET("/objects/:file", ObjectsHandler)
	e.Static("/static", "public")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	e.Logger.Fatal(e.Start(":1323"))
}
