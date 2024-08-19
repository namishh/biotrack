package services

import (
	"fmt"

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
	stmt := `INSERT INTO profile (profile_of, profile_picture) VALUES (?, ?)`
	avatarLink := fmt.Sprintf("/avatar/%s", u.Username)
	_, err := ps.ProfileStore.DB.Exec(stmt, u.ID, avatarLink)

	return err
}

func (ps *ProfileService) GetProfileByUserId(id int) (Profile, error) {
	profile := Profile{}
	stmt := `SELECT id, level, profile_picture, weight, height, birthday, bio, profile_of, last_login FROM profiles WHERE profile_of = ?`

	row := ps.ProfileStore.DB.QueryRow(stmt, id)

	err := row.Scan(&profile.ID, &profile.Level, &profile.ProfilePicture, &profile.Weight, &profile.Height, &profile.Birthday, &profile.Bio, &profile.ProfileOf, &profile.LastLogin)

	if err != nil {
		return Profile{}, err
	}

	return profile, nil
}
