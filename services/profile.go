package services

import (
	"fmt"
	"log"

	"github.com/namishh/biotrack/database"
)

type Profile struct {
	ID             int     `json:"id"`
	Level          int     `json:"level"`
	ProfilePicture string  `json:"profile_picture"`
	Weight         float64 `json:"weight"`
	Height         float64 `json:"height"`
	Birthday       string  `json:"birthday"`
	Streak         int     `json:"streak"`
	Bio            string  `json:"bio"`
	ProfileOf      int     `json:"profile_of"`
	LastLogin      string  `json:"last_login"`
}

type ProfileService struct {
	Profile      Profile
	ProfileStore database.DatabaseStore
}

func NewProfileService(profile Profile, profileStore database.DatabaseStore) *ProfileService {
	return &ProfileService{
		Profile:      profile,
		ProfileStore: profileStore,
	}
}

func (ps *ProfileService) CreateDefaultProfile(u User) error {
	// create user himself
	stmt := `INSERT INTO profiles (profile_of, profile_picture) VALUES ($1, $2)`
	log.Print(u)
	avatarLink := fmt.Sprintf("/avatar/%s", u.Username)
	_, err := ps.ProfileStore.DB.Exec(stmt, u.ID, avatarLink) 

	return err
}
