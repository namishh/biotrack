package services

import (
	"github.com/namishh/biotrack/database"
)

type Chat struct {
	Id        int    `json:"id"`
	Message   string `json:"message"`
	Sender    string `json:"sender"`
	CreatedAt string `json:"created_at"`
}

type Chart struct {
	Id        int    `json:"id"`
	Labels    string `json:"labels"`
	Data      string `json:"data"`
	ChatId    int    `json:"chat_id"`
	CreatedAt string `json:"created_at"`
}

type ChatService struct {
	Chat      Chat
	Chart     Chart
	ChatStore database.DatabaseStore
}

func NewChatService(chat Chat, chart Chart, store database.DatabaseStore) *ChatService {
	return &ChatService{Chat: chat, Chart: chart, ChatStore: store}
}
