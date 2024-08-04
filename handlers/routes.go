package handlers

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

// Handler
func hello(c echo.Context) error {
	return c.Render(http.StatusOK, "index", "World")
}
