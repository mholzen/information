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
	"github.com/mholzen/information/triples/transforms"
	"github.com/mholzen/information/triples/transforms/node_link"

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
	tm := transforms.NewFileJsonParser(path)
	res := triples.NewTriples()
	err := res.Transform(tm.Transformer)
	if err != nil {
		return nil, echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	return res, nil
}

func HtmlHandler(c echo.Context) error {
	src, err := GetTripleList(c)
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
	src, err := GetTripleList(c)
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
	src, err := GetTripleList(c)
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
	src, err := GetTripleList(c)
	if err != nil {
		return err
	}
	def := transforms.NewDefaultTableDefinition(src)
	tr := transforms.NewTableGenerator(def)
	err = src.Transform(tr.Transformer)
	if err != nil {
		return err
	}
	return c.HTML(http.StatusOK, tr.Html())
}

func main() {
	e := echo.New()
	e.GET("/html/:file", HtmlHandler)
	e.GET("/objects/:file", ObjectsHandler)
	e.GET("/nodelink/:file", NodeLinkHandler)
	e.GET("/graph/:file", GraphHandler)
	e.GET("/table/:file", TableHandler)

	e.Static("/static", "public")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	e.Logger.Fatal(e.Start(":1323"))
}
