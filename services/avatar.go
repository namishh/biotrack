package services

import (
	"math"

	"github.com/lucasb-eyer/go-colorful"
	"github.com/namishh/biotrack/database"
)

type Avatar struct {
	Id        int
	Username  string
	FromColor string
	ToColor   string
}

type AvatarService struct {
	Avatar      Avatar
	AvatarStore database.DatabaseStore
}

func NewAvatarService(avatar Avatar, avatarStore database.DatabaseStore) *AvatarService {
	return &AvatarService{Avatar: avatar, AvatarStore: avatarStore}
}

func djb2(str string) int {
	hash := 5381
	for _, char := range str {
		hash = ((hash << 5) + hash) + int(char)
	}
	return hash
}

func (as *AvatarService) GenerateGradient(username string) map[string]string {
	h := float64(djb2(username)%360) / 360.0 // Normalize to [0, 1]
	c1 := colorful.Hsv(h*360, 0.95, 1)

	// Calculate the triad color
	h2 := math.Mod(h+1.0/3, 1.0) // Add 1/3 for triad, keep in [0, 1]
	c2 := colorful.Hsv(h2*360, 0.95, 1)

	return map[string]string{
		"fromColor": c1.Hex(),
		"toColor":   c2.Hex(),
	}
}
