package main

import (
	"context"
	"io"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"text/template"
	"time"

	"github.com/mholzen/information/handlers"
	"github.com/mholzen/information/triples/data"

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

	err := data.InitData()
	if err != nil {
		e.Logger.Fatal(err)
	}

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

	e.Static("/static", "public")

	e.GET("/stats/:file", handlers.StatsHandler)
	e.GET("/triples/:file", handlers.TriplesHandler)
	e.GET("/html/:file", handlers.HtmlHandler)
	e.GET("/nodelink/:file", handlers.NodeLinkHandler)
	e.GET("/graph/:file", handlers.GraphHandler)
	e.GET("/files/:file", handlers.FilesPostfixHandler)
	e.GET("/data", handlers.FilesPostfixHandler)

	e.GET("/", func(c echo.Context) error {
		return c.Redirect(http.StatusFound, "/files/data/index.md/text/html")
	})

	renderer := &TemplateRenderer{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.Renderer = renderer

	go func() {
		if err := e.Start(":1323"); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server with
	// a timeout of 10 seconds.
	quit := make(chan os.Signal, 1)

	// signal.Notify registers the given channel to receive notifications of the specified signals.
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	e.Logger.Info("Signal received: shutting down the server")

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
