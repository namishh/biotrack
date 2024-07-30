package main

import (
	"html/template"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	t := &Template{
		templates: template.Must(template.ParseGlob("views/*/*.html")),
	}

	//Css styles
	e.Static("/css", "public/assets")

	e.Renderer = t
	e.GET("/", hello)

	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}

// Handler
func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}
