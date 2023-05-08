package main

import (
	"io"
	"net/http"
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
	all, err := triples.Parse(c.Param("file"))
	if err != nil {
		return err
	}
	tripleList := all.GetTripleList()
	// tripleList = make(triples.TripleList, 0)
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("a"), triples.NewStringNode("b"), triples.NewStringNode("c")})
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("a"), triples.NewStringNode("b"), triples.NewStringNode("d")})
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("a"), triples.NewStringNode("b1"), triples.NewStringNode("c")})
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("a"), triples.NewStringNode("b2"), triples.NewStringNode("c")})
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("b"), triples.NewStringNode("b2"), triples.NewStringNode("c")})
	// tripleList = append(tripleList, triples.Triple{triples.NewStringNode("zz"), triples.NewStringNode("b"), triples.NewStringNode("c")})
	tripleList.Sort()

	// subject := triples.NewStringNode("marc")
	// marc := all.AddReachableTriples(subject, nil)
	// subjectTriples := marc.GetTriplesForSubject(subject, nil)
	html := triples.NewHtmlTransformer(*all, tripleList, 0)

	data := map[string]interface{}{
		"tripleList": tripleList,
		"html":       html.String(),
	}
	return c.Render(http.StatusOK, "index.html", data)
}

func init() {
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
