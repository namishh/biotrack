package main

import (
	"log"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/namishh/biotrack/database"
	"github.com/namishh/biotrack/handlers"
	"github.com/namishh/biotrack/services"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	// Echo instance
	e := echo.New()

	SECRET_KEY := os.Getenv("SECRET")
	DB_NAME := os.Getenv("DB_NAME")

	// Use Middleware Here
	e.Use(middleware.Logger())

	e.Use(session.Middleware(sessions.NewCookieStore([]byte(SECRET_KEY))))

	e.Static("/static", "public/assets")

	store, err := database.NewDatabaseStore(DB_NAME)

	us := services.NewUserService(services.User{}, store)
	ah := handlers.NewAuthHandler(us)

	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	handlers.SetupRoutes(e, ah)

	// Start server
	e.Logger.Fatal(e.Start(":6969"))
}
