package handlers

import (
	"errors"
	"net/http"
	"net/mail"

	"github.com/a-h/templ"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages"
)

const auth_key = "auth_key"
const auth_sessions_key = "auth_session_key"
const user_id_key = "user_id_key"
const user_name_key = "user_name_key"
const tzone_key = "tzone_key"

type AuthService interface {
	CreateUser(u services.User) error
	CheckEmail(email string) (services.User, error)
}

type AuthHandler struct {
	UserServices AuthService
}

func NewAuthHandler(us AuthService) *AuthHandler {
	return &AuthHandler{UserServices: us}
}

func renderView(c echo.Context, cmp templ.Component) error {
	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return cmp.Render(c.Request().Context(), c.Response().Writer)
}

func (ah *AuthHandler) flagsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			c.Set("FROMPROTECTED", false)

			return next(c)
		}

		c.Set("FROMPROTECTED", true)

		return next(c)
	}
}

func (ah *AuthHandler) authMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		sess, _ := session.Get(auth_sessions_key, c)
		if auth, ok := sess.Values[auth_key].(bool); !ok || !auth {
			c.Set("FROMPROTECTED", false)
			return echo.NewHTTPError(echo.ErrUnauthorized.Code, "Please provide valid credentials")
		}

		if userId, ok := sess.Values[user_id_key].(int); ok && userId != 0 {
			c.Set(user_id_key, userId) // set the user_id in the context
		}

		if username, ok := sess.Values[user_name_key].(string); ok && len(username) != 0 {
			c.Set(user_name_key, username) // set the username in the context
		}

		if tzone, ok := sess.Values[tzone_key].(string); ok && len(tzone) != 0 {
			c.Set(tzone_key, tzone) // set the client's time zone in the context
		}

		// fromProtected = true
		c.Set("FROMPROTECTED", true)

		return next(c)
	}
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
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	// isError = false
	view := pages.Login(fromProtected, errs)
	c.Set("ISERROR", false)

	return renderView(c, pages.LoginIndex(
		"Login",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}

func (ah *AuthHandler) RegisterHandler(c echo.Context) error {

	errs := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if c.Request().Method == "POST" {
		// Validating data here
		email := c.FormValue("email")
		password := c.FormValue("password")
		username := c.FormValue("username")

		// check if email is valid
		if !valid(email) {
			errs["email"] = "Invalid email address"
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

		if errs["username"] != "" || errs["email"] != "" || errs["password"] != "" {
			view := pages.Register(fromProtected, errs)

			c.Set("ISERROR", false)

			return renderView(c, pages.RegisterIndex(
				"Register",
				"",
				fromProtected,
				c.Get("ISERROR").(bool),
				view,
			))
		}

		return c.Redirect(http.StatusSeeOther, "/login")
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	view := pages.Register(fromProtected, errs)

	c.Set("ISERROR", false)

	return renderView(c, pages.RegisterIndex(
		"Register",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}
