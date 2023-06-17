package main

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
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

func Handler(c echo.Context) error {
	root := os.Getenv("ROOT")
	path := filepath.Join(root, c.Param("file"))
	all, err := triples.Parse(path)
	if err != nil {
		return echo.NewHTTPError(http.StatusNotFound, err.Error())
	}
	tripleList := all.GetTripleList()
	tripleList.Sort()

	html := triples.NewHtmlTransformer(*all, tripleList, 0)

	data := map[string]interface{}{
		"tripleList": tripleList,
		"html":       html.String(),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func main() {
	e := echo.New()
	e.GET("/html/:file", Handler)
	e.Static("/static", "public")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	e.Logger.Fatal(e.Start(":1323"))
}
