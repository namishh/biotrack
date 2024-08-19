package handlers

import (
	"errors"
	"log"
	"net/http"
	"net/mail"
	"strings"

	"github.com/a-h/templ"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages"
	"github.com/namishh/biotrack/views/pages/auth"
	"github.com/namishh/biotrack/views/pages/profile"
	"golang.org/x/crypto/bcrypt"
)

const auth_key string = "auth_key"
const auth_sessions_key string = "auth_session_key"
const user_id_key string = "user_id_key"
const user_name_key string = "user_name_key"
const tzone_key string = "tzone_key"

type AuthService interface {
	CreateUser(u services.User) (services.User, error)
	CheckEmail(email string) (services.User, error)
	CheckUsername(usr string) (services.User, error)
}

type ProfileService interface {
	CreateDefaultProfile(u services.User) error
	GetProfileByUserId(id int) (services.Profile, error)
}

type AvatarService interface {
	GenerateGradient(username string) map[string]string
}

type AuthHandler struct {
	UserServices    AuthService
	ProfileServices ProfileService
	AvatarServices  AvatarService
}

func NewAuthHandler(us AuthService, ps ProfileService, as AvatarService) *AuthHandler {
	return &AuthHandler{UserServices: us, ProfileServices: ps, AvatarServices: as}
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func valid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func (ah *AuthHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	// isError = false
	homeView := pages.Home(fromProtected)
	c.Set("ISERROR", false)

	return renderView(c, pages.HomeIndex(
		"Home",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		homeView,
	))
}

func (ah *AuthHandler) LoginHandler(c echo.Context) error {
	errs := make(map[string]string)
	formdata := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	if c.Request().Method == "POST" {
		tzone := ""
		if len(c.Request().Header["X-Timezone"]) != 0 {
			tzone = c.Request().Header["X-Timezone"][0]
		}

		formdata["email"] = c.FormValue("email")
		formdata["password"] = c.FormValue("password")

		log.Print(tzone)

		user, err := ah.UserServices.CheckEmail(c.FormValue("email"))
		if err != nil {
			if strings.Contains(err.Error(), "no rows in result set") {
				c.Set("ISERROR", true)
				errs["dne"] = "User with this email does not exist."
				view := auth.Login(fromProtected, formdata, errs)

				return renderView(c, auth.LoginIndex(
					"Login",
					"",
					fromProtected,
					c.Get("ISERROR").(bool),
					view,
				))
			}
		}

		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(c.FormValue("password")),
		)
		if err != nil {
			c.Set("ISERROR", true)
			errs["pass"] = "Incorrect Password"
			view := auth.Login(fromProtected, formdata, errs)

			return renderView(c, auth.LoginIndex(
				"Login",
				"",
				fromProtected,
				c.Get("ISERROR").(bool),
				view,
			))
		}

		// Log in the user
		sess, _ := session.Get(auth_sessions_key, c)
		sess.Options = &sessions.Options{
			Path:     "/",
			MaxAge:   60 * 60 * 24 * 7, // 1 week
			HttpOnly: true,
		}

		// Set user as authenticated, their username,
		// their ID and the client's time zone
		sess.Values = map[interface{}]interface{}{
			auth_key:      true,
			user_id_key:   user.ID,
			user_name_key: user.Username,
			tzone_key:     tzone,
		}
		sess.Save(c.Request(), c.Response())

		return c.Redirect(http.StatusSeeOther, "/")

	}
	// isError = false
	view := auth.Login(fromProtected, formdata, errs)
	c.Set("ISERROR", false)

	return renderView(c, auth.LoginIndex(
		"Login",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}

func (ah *AuthHandler) RegisterHandler(c echo.Context) error {

	errs := make(map[string]string)
	formdata := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request().Method == "POST" {
		// Validating data here
		email := c.FormValue("email")
		formdata["email"] = email
		password := c.FormValue("password")
		formdata["password"] = password
		username := c.FormValue("username")
		formdata["username"] = username

		// check if email is valid
		if !valid(email) {
			errs["email"] = "Invalid email address"
			c.Set("ISERROR", true)
		}

		_, err := ah.UserServices.CheckEmail(email)
		if err == nil {
			errs["email"] = "Account with this email already exists"
			c.Set("ISERROR", true)
		}

		// password valid: minimum 4 letters
		if len(password) < 4 {
			errs["password"] = "Password must be at least 4 characters"
		}

		// username valid: minimum 4 letters
		if len(username) < 4 {
			errs["username"] = "Username must be at least 4 characters"
		}

		_, err = ah.UserServices.CheckUsername(username)
		log.Print(err)
		if err == nil {
			errs["username"] = "Account with this username already exists"
			c.Set("ISERROR", true)
		}

		if errs["username"] != "" || errs["email"] != "" || errs["password"] != "" {
			view := auth.Register(fromProtected, formdata, errs)

			c.Set("ISERROR", false)

			return renderView(c, auth.RegisterIndex(
				"Register",
				"",
				fromProtected,
				c.Get("ISERROR").(bool),
				view,
			))
		}

		user := services.User{
			Email:    email,
			Username: username,
			Password: password,
		}

		u, err := ah.UserServices.CreateUser(user)
		if err != nil {
			return err
		}
		ah.ProfileServices.CreateDefaultProfile(u)

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	view := auth.Register(fromProtected, formdata, errs)

	c.Set("ISERROR", false)

	return renderView(c, auth.RegisterIndex(
		"Register",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}

func (ah *AuthHandler) LogoutHandler(c echo.Context) error {
	sess, _ := session.Get(auth_sessions_key, c)
	fromProtected, _ := c.Get("FROMPROTECTED").(bool)

	if !fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}
	// Revoke users authentication
	sess.Values = map[interface{}]interface{}{
		auth_key:      false,
		user_id_key:   "",
		user_name_key: "",
		tzone_key:     "",
	}
	sess.Save(c.Request(), c.Response())

	// fromProtected = false
	c.Set("FROMPROTECTED", false)

	return c.Redirect(http.StatusSeeOther, "/login")
}

func (ah *AuthHandler) ProfileHandler(c echo.Context) error {
	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if !fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if c.Request().Method == "POST" {
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	log.Print(c.Get(user_id_key))
	user, err := ah.UserServices.CheckUsername(c.Get(user_name_key).(string))

	if err != nil {
		return c.Redirect(200, "/login")
	}

	p := services.Profile{}

	view := profile.Profile(fromProtected, user, p, errs)

	c.Set("ISERROR", false)

	return renderView(c, profile.ProfileIndex(
		"Profile",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}
