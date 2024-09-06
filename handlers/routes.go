package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ah *AuthHandler, jh *JournalHandler, ch *ChatHandler) {
	e.GET("/", flagsMiddleware(ah.HomeHandler))

	// AUTH ROUTES
	e.GET("/register", flagsMiddleware(ah.RegisterHandler))
	e.POST("/register", flagsMiddleware(ah.RegisterHandler))

	e.GET("/login", flagsMiddleware(ah.LoginHandler))
	e.POST("/login", flagsMiddleware(ah.LoginHandler))

	e.GET("/logout", flagsMiddleware(ah.LogoutHandler))
	e.GET("/avatar/:name", flagsMiddleware(ah.Avatar))

	e.GET("/profile", authMiddleware(ah.ProfileHandler))
	e.POST("/profile", authMiddleware(ah.ProfileHandler))

	journalGroup := e.Group("/journal", authMiddleware)
	journalGroup.GET("", jh.HomeHandler)
	journalGroup.GET("/calendar", jh.CalendarHandler)
	journalGroup.GET("/new", jh.NewHandler)
	journalGroup.GET("/:year/:month", jh.MonthHandler)
	journalGroup.GET("/:year/:month/:date", jh.DayHandler)
	journalGroup.GET("/:year/:month/:date/delete/:id", jh.DeleteHandler)
	journalGroup.POST("/:year/:month/:date", jh.DayHandler)

	chatGroup := e.Group("/chat", authMiddleware)
	chatGroup.GET("", ch.HomeHandler)
	chatGroup.POST("", ch.HomeHandler)

	e.GET("/*", RouteNotFoundHandler)
}
