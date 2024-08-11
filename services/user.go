package services

import (
	"database/sql"
	"fmt"

	"github.com/namishh/biotrack/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
	Profile   Profile
}

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

type UserService struct {
	User      User
	UserStore database.DatabaseStore
}

func NewUserService(user User, userStore database.DatabaseStore) *UserService {
	return &UserService{
		User:      user,
		UserStore: userStore,
	}
}

func (us *UserService) CreateUser(u User) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// create user himself
	stmt := `INSERT INTO users (email, password, username) VALUES ($1, $2, $3)`

	_, err = us.UserStore.DB.Exec(stmt, u.Email, string(hashedPassword), u.Username)

	return err
}

func (us *UserService) CreateDefaultProfile(u User) (sql.Result, error) {
	// create user himself
	stmt := `INSERT INTO profiles (profile_of, profile_picture) VALUES ($1, $2)`

	avatarLink := fmt.Sprintf("/avatar/%s", u.Username)
	profile, err := us.UserStore.DB.Exec(stmt, u.ID, avatarLink) // sorry vercel dont be mad

	return profile, err
}

func (us *UserService) CheckEmail(email string) (User, error) {
	query := `SELECT id, email, password, username FROM users
		WHERE email = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Email = email
	err = stmt.QueryRow(
		us.User.Email,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}

func (us *UserService) CheckProfile(id int) (Profile, error) {
	query := `SELECT id, level, profile_picture, weight, height, birthday, streak, bio, profile_of, last_login FROM profiles
		WHERE profile_of = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return Profile{}, err
	}

	defer stmt.Close()

	us.User.ID = id
	err = stmt.QueryRow(
		us.User.ID,
	).Scan(
		&us.User.Profile.ID,
		&us.User.Profile.Level,
		&us.User.Profile.ProfilePicture,
		&us.User.Profile.Weight,
		&us.User.Profile.Height,
		&us.User.Profile.Birthday,
		&us.User.Profile.Streak,
		&us.User.Profile.Bio,
		&us.User.Profile.ProfileOf,
		&us.User.Profile.LastLogin,
	)
	if err != nil {
		return Profile{}, err
	}

	return us.User.Profile, nil
}

func (us *UserService) CheckUsername(usr string) (User, error) {
	query := `SELECT id, email, password, username FROM users
		WHERE username = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.Username = usr
	err = stmt.QueryRow(
		us.User.Username,
	).Scan(
		&us.User.ID,
		&us.User.Email,
		&us.User.Password,
		&us.User.Username,
	)
	if err != nil {
		return User{}, err
	}

	return us.User, nil
}
