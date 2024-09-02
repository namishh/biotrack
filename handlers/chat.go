package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/views/pages/chat"
)

type ChatService interface{}

type ChatHandler struct {
	ChatService ChatService
}

func NewChatHandler(cs ChatService) *ChatHandler {
	return &ChatHandler{ChatService: cs}
}

func (ch *ChatHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	cview := chat.Home(fromProtected)
	c.Set("ISERROR", false)

	return renderView(c, chat.HomeIndex(
		"Chat",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		cview,
	))
}
