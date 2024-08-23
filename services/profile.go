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

func (ps *ProfileService) UpdateProfile(userid int, height float64, weight float64, dob string) error {

	stmt := `UPDATE profile SET weight = ?, height = ?, birthday = ? WHERE profile_of = ?`
	_, err := ps.ProfileStore.DB.Exec(stmt, weight, height, dob, userid)

	return err
}

func (ps *ProfileService) UpdateProfilePicture(u User, pfp string) error {
	stmt := `UPDATE profile SET profile_picture = ? WHERE profile_of = ?`
	_, err := ps.ProfileStore.DB.Exec(stmt, pfp, u.ID)

	return err
}

func (ps *ProfileService) GetProfileByUserId(id int) (Profile, error) {
	query := `SELECT id, level, profile_picture, weight, height, birthday, bio FROM profile WHERE profile_of = ?`

	stmt, err := ps.ProfileStore.DB.Prepare(query)
	if err != nil {
		log.Print(err)
		return Profile{}, err
	}

	defer stmt.Close()

	ps.Profile.ProfileOf = id

	err = stmt.QueryRow(ps.Profile.ProfileOf).Scan(&ps.Profile.ID, &ps.Profile.Level, &ps.Profile.ProfilePicture, &ps.Profile.Weight, &ps.Profile.Height, &ps.Profile.Birthday, &ps.Profile.Bio)

	if err != nil {
		log.Print(err)
		return Profile{}, err
	}

	return ps.Profile, nil
}
