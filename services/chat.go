package services

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"github.com/namishh/biotrack/database"
	"google.golang.org/api/option"
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
	AI        *genai.ChatSession
}

func NewChatService(chat Chat, chart Chart, store database.DatabaseStore) (*ChatService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI_API_KEY")))
	if err != nil {
		return nil, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")
	aiSession := model.StartChat()

	// Load background from file
	prompt, err := os.ReadFile("./assets/Prompt.txt")
	if err != nil {
		return nil, err
	}

	aiSession.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(prompt),
			},
			Role: "user",
		},
	}

	return &ChatService{
		Chat:      chat,
		Chart:     chart,
		ChatStore: store,
		AI:        aiSession,
	}, nil
}

func (cs *ChatService) NewUserChat(userid int, message string) error {
	stmt := `INSERT INTO chat (sender, message) VALUES (?, ?)`
	_, err := cs.ChatStore.DB.Exec(stmt, userid, message)

	return err
}

func (cs *ChatService) NewAIChat(userid int, message string) error {
	stmt := `INSERT INTO chat (sender, message) VALUES (?, ?)`
	_, err := cs.ChatStore.DB.Exec(stmt, fmt.Sprintf("AI-%d", userid), message)

	return err
}

func (cs *ChatService) GenerateResponse() error {
	return nil
}
