package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/views/pages/profile"
)

func (ah *AuthHandler) ProfileHandler(c echo.Context) error {
	errs := make(map[string]string)
	formdata := make(map[string]string)
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if !fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	p, err := ah.ProfileServices.GetProfileByUserId(c.Get(user_id_key).(int))

	user, err := ah.UserServices.CheckUsername(c.Get(user_name_key).(string))

	if err != nil {
		return c.Redirect(200, "/login")
	}

	view := profile.Profile(fromProtected, user, p, errs, formdata)

	c.Set("ISERROR", false)

	return renderView(c, profile.ProfileIndex(
		"Profile",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		view,
	))
}

func (ah *AuthHandler) UpdateProfileAvatarHandler(c echo.Context) error {
	if c.Request().Method == "GET" {
		return c.Redirect(http.StatusSeeOther, "/profile")
	}
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)

	if !fromProtected {
		return c.Redirect(http.StatusSeeOther, "/")
	}

	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	form, err := c.MultipartForm()
	if err != nil {
		return c.Redirect(200, "/profile")
	}

	files := form.File["profile_picture"]

	if len(files) == 0 {
		return c.Redirect(200, "/profile")
	}

	av := files[0]

	src, err := av.Open()
	if err != nil {
		return c.Redirect(200, "/profile")
	}

	defer src.Close()

	filename := fmt.Sprintf("./public/avatars/%s.img", c.Get(user_name_key).(string))
	finalurl := fmt.Sprintf("/av/%s.img", c.Get(user_name_key).(string))


	dst, err := os.Create(filename)

	if err != nil {
		return c.Redirect(200, "/profile")
	}

	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.Redirect(200, "/profile")
	}

	_, err = ah.ProfileServices.GetProfileByUserId(c.Get(user_id_key).(int))

	user, err := ah.UserServices.CheckUsername(c.Get(user_name_key).(string))

	ah.ProfileServices.UpdateProfilePicture(user, finalurl)

	if err != nil {
		return c.Redirect(200, "/login")
	}

	return c.Redirect(200, "/profile")
}
