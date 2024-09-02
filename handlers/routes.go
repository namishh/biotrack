package handlers

import (
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, ah *AuthHandler, jh *JournalHandler, ch *ChatHandler) {
	e.GET("/", ah.flagsMiddleware(ah.HomeHandler))

	// AUTH ROUTES
	e.GET("/register", ah.flagsMiddleware(ah.RegisterHandler))
	e.POST("/register", ah.flagsMiddleware(ah.RegisterHandler))

	e.GET("/login", ah.flagsMiddleware(ah.LoginHandler))
	e.POST("/login", ah.flagsMiddleware(ah.LoginHandler))

	e.GET("/logout", ah.flagsMiddleware(ah.LogoutHandler))
	e.GET("/avatar/:name", ah.flagsMiddleware(ah.Avatar))

	e.GET("/profile", ah.authMiddleware(ah.ProfileHandler))
	e.POST("/profile", ah.authMiddleware(ah.ProfileHandler))

	journalGroup := e.Group("/journal", ah.authMiddleware)
	journalGroup.GET("", jh.HomeHandler)
	journalGroup.GET("/calendar", jh.CalendarHandler)
	journalGroup.GET("/new", jh.NewHandler)
	journalGroup.GET("/:year/:month", jh.MonthHandler)
	journalGroup.GET("/:year/:month/:date", jh.DayHandler)
	journalGroup.GET("/:year/:month/:date/delete/:id", jh.DeleteHandler)
	journalGroup.POST("/:year/:month/:date", jh.DayHandler)
	
	chatGroup := e.Group("/chat", ah.authMiddleware)
	chatGroup.GET("", ch.HomeHandler)
}
