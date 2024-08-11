package handlers

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
)

const size = 192

func createSVG(gradient map[string]string) string {
	return fmt.Sprintf(`<svg width="%d" height="%d" viewBox="0 0 %d %d" version="1.1" xmlns="http://www.w3.org/2000/svg">
		<defs>
			<linearGradient id="gradient" x1="0" y1="0" x2="1" y2="1">
				<stop offset="0%%" stop-color="%s" />
				<stop offset="100%%" stop-color="%s" />
			</linearGradient>
		</defs>
		<rect fill="url(#gradient)" x="0" y="0" width="%d" height="%d" />
	</svg>`, size, size, size, size, gradient["fromColor"], gradient["toColor"], size, size)
}

func Avatar(c echo.Context) error {
	name := c.Param("user")
	parts := strings.Split(name, ".")
	username := parts[0]

	if username == "" {
		username = fmt.Sprintf("%f", rand.Float64())
	}

	gradient := services.GenerateGradient(username)
	log.Print(gradient)

	svg := createSVG(gradient)
	c.Response().Header().Set("Cache-Control", "public, max-age=604800, immutable")
	return c.Blob(http.StatusOK, "image/svg+xml", []byte(svg))
	// For PNG, we'd need to implement image generation
	// This is a placeholder and won't actually generate a PNG
}
