package services

import (
	"github.com/namishh/biotrack/database"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID        int    `json:"id"`
	Email     string `json:"email"`
	Password  string `json:"password"`
	Username  string `json:"username"`
	CreatedAt string `json:"created_at"`
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
	User         User
	Profile      Profile
	ProfileStore database.DatabaseStore
	UserStore    database.DatabaseStore
}

func NewUserService(user User, userStore database.DatabaseStore, profile Profile, profileStore database.DatabaseStore) *UserService {
	return &UserService{
		User:         user,
		UserStore:    userStore,
		Profile:      profile,
		ProfileStore: profileStore,
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
	if err != nil {
		return err
	}

	// create user profile
	stmt = `INSERT INTO profiles (profile_of, profile_picture) VALUES ($1, $2)`
	_, err = us.ProfileStore.DB.Exec(stmt, u.ID, "default.jpg")

	return err
}

func (us *UserService) GetUserByEmail(email string) (User, error) {
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
