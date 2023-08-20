package main

import (
	"io"
	"net/http"
	"text/template"

	"github.com/mholzen/information/handlers"

	"github.com/labstack/echo/v4"
)

type TemplateRenderer struct {
	templates *template.Template
}

func (t *TemplateRenderer) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func main() {
	e := echo.New()

	e.HTTPErrorHandler = func(err error, c echo.Context) {
		code := http.StatusInternalServerError
		if he, ok := err.(*echo.HTTPError); ok {
			code = he.Code
		}
		c.JSON(code, map[string]interface{}{
			"status":  code,
			"message": err.Error(),
		})
	}

	e.GET("/stats/:file", handlers.StatsHandler)
	e.GET("/triples/:file", handlers.TriplesHandler)
	e.GET("/html/:file", handlers.HtmlHandler)
	e.GET("/objects/:file", handlers.ObjectsHandler)
	e.GET("/nodelink/:file", handlers.NodeLinkHandler)
	e.GET("/graph/:file", handlers.GraphHandler)
	e.GET("/table/:file", handlers.TableHandler)

	e.Static("/static", "public")
	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer
	e.Logger.Fatal(e.Start(":1323"))
}
