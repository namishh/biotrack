package handlers

import (
	"errors"
	"fmt"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages/chat"
)

type ChatService interface {
	NewUserChat(userid int, message string) error
	GetAllChatsByUserId(userid int) ([]services.Chat, error)
	GenerateResponse(id int, ud string) error
	NewAIChat(userid int, message string) error
}

type ChatHandler struct {
	ChatService  ChatService
	EntryService EntryService
}

func NewChatHandler(cs ChatService, es EntryService) *ChatHandler {
	return &ChatHandler{ChatService: cs, EntryService: es}
}

func (ch *ChatHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		if len(c.FormValue("message")) > 0 {
		ch.ChatService.NewUserChat(c.Get(user_id_key).(int), c.FormValue("message"))
		fd, err := ch.EntryService.GetFormattedEntriesByUser(c.Get(user_id_key).(int))
		err = ch.ChatService.GenerateResponse(c.Get(user_id_key).(int), fd)
		fmt.Println(err)
		}
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
