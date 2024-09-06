package services

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

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
	Client    *genai.Client
	Context   context.Context
	Chat      Chat
	Chart     Chart
	ChatStore database.DatabaseStore
	AI        *genai.ChatSession
}

func NewChatService(chat Chat, chart Chart, store database.DatabaseStore) (*ChatService, error) {
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("GEMINI")))
	if err != nil {
		return nil, err
	}

	model := client.GenerativeModel("gemini-1.5-flash")
	aiSession := model.StartChat()

	log.Println("Chat service created")
	return &ChatService{
		Client:    client,
		Context:   ctx,
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

func (cs *ChatService) GetAllChatsByUserId(userid int) ([]Chat, error) {
	query := `
        SELECT id, sender, message, created_at
        FROM chat
        WHERE sender = ? OR sender = ?
        ORDER BY id ASC
    `

	rows, err := cs.ChatStore.DB.Query(query, userid, fmt.Sprintf("AI-%d", userid))
	if err != nil {
		return nil, fmt.Errorf("error querying chats: %v", err)
	}
	defer rows.Close()

	var chats []Chat

	for rows.Next() {
		var chat Chat
		err := rows.Scan(&chat.Id, &chat.Sender, &chat.Message, &chat.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("error scanning chat row: %v", err)
		}
		chats = append(chats, chat)
	}

	// Check for errors from iterating over rows
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating chat rows: %v", err)
	}

	return chats, nil
}

func splitChats(chats []Chat) ([]Chat, Chat) {
	if len(chats) == 0 {
		return nil, Chat{}
	}

	if len(chats) == 1 {
		return nil, chats[0]
	}

	allButLast := chats[:len(chats)-1]
	last := chats[len(chats)-1]

	return allButLast, last
}

func convertPartToString(part genai.Part) string {
	switch v := part.(type) {
	case genai.Text:
		return string(v)
	default:
		return "NIL"
	}
}

func (cs *ChatService) GenerateResponse(id int, ud string) error {
	if cs.AI == nil {
		return fmt.Errorf("AI session is not initialized")
	}
	log.Print(ud)
	chats, err := cs.GetAllChatsByUserId(id)
	if err != nil {
		return err
	}

	prompt, err := os.ReadFile("./assets/PROMPT.txt")
	if err != nil {
		return err
	}

	cs.AI.History = []*genai.Content{
		{
			Parts: []genai.Part{
				genai.Text(prompt),
			},
			Role: "user",
		},
		{
			Parts: []genai.Part{
				genai.Text(ud),
			},
			Role: "user",
		},
	}
	chs, last := splitChats(chats)
	for _, chat := range chs {
		_, err := strconv.Atoi(chat.Sender)
		sender := "user"
		if err != nil {
			sender = "model"
		}
		cs.AI.History = append(cs.AI.History, &genai.Content{
			Parts: []genai.Part{
				genai.Text(chat.Message),
			},
			Role: sender,
		})
	}

	res, err := cs.AI.SendMessage(cs.Context, genai.Text(last.Message))
	log.Print(err)
	if err != nil {
		return err
	}

	str := ""
	for _, cand := range res.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				str = str + convertPartToString(part)
			}
		}
	}
	log.Print(str)
	err = cs.NewAIChat(id, str)

	return err
}
