package services

import (
	"log"

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

func (us *UserService) CreateUser(u User) (User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	// create user himself
	stmt := `INSERT INTO users (email, password, username) VALUES (?, ?, ?) RETURNING id`
	err = us.UserStore.DB.QueryRow(stmt, u.Email, string(hashedPassword), u.Username).Scan(&u.ID)

	return u, err
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

func (us *UserService) UpdateUser(email string, username string, id int) error {
	stmt := `UPDATE users SET email = ?, username = ? WHERE id = ?`
	_, err := us.UserStore.DB.Exec(stmt, email, username, id)
	log.Println("done")

	if err != nil {
		return err
	}

	return nil
}


func (us *UserService) UpdateUsername(username string, id int) error {
	stmt := `UPDATE users SET username = ? WHERE id = ?`
	_, err := us.UserStore.DB.Exec(stmt, username, id)

	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UpdateEmail(email string, id int) error {
	stmt := `UPDATE users SET email = ? WHERE id = ?`
	_, err := us.UserStore.DB.Exec(stmt, email, id)

	if err != nil {
		return err
	}

	return nil
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

func (us *UserService) CheckID(usr int) (User, error) {
	query := `SELECT id, email, password, username FROM users
		WHERE id = ?`

	stmt, err := us.UserStore.DB.Prepare(query)
	if err != nil {
		return User{}, err
	}

	defer stmt.Close()

	us.User.ID = usr
	err = stmt.QueryRow(
		us.User.ID,
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