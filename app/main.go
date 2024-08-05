package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// Echo instance
	e := echo.New()

	// Use Middleware Here
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "public/assets")

	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}
