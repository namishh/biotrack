package handlers

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages/profile"
)

type ProfileService interface {
	CreateDefaultProfile(u services.User) error
	UpdateProfilePicture(u services.User, pfp string) error
	UpdateProfile(userid int, height float64, weight float64, dob string, heightunit string, weightunit string) error
	GetProfileByUserId(id int) (services.Profile, error)
	UpdateProfileWeight(u int, weight int) error
	UpdateProfileHeight(u int, height int) error
}

const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

func StringWithCharset(length int, charset string) string {
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func String(length int) string {
	return StringWithCharset(length, charset)
}

func ValidateDate(dateString string) bool {
	pattern := `^\d{4}-(0[1-9]|1[0-2])-(0[1-9]|[12][0-9]|3[01])$`
	matched, err := regexp.MatchString(pattern, dateString)
	if err != nil || !matched {
		return false
	}

	birthDate, err := time.Parse("2006-01-02", dateString)
	if err != nil {
		return false
	}

	if birthDate.After(time.Now()) {
		return false
	}

	today := time.Now()
	age := today.Year() - birthDate.Year()

	if today.Month() < birthDate.Month() || (today.Month() == birthDate.Month() && today.Day() < birthDate.Day()) {
		age--
	}

	return age >= 13
}

func RoundFloat(val float64, precision uint) float64 {
	ratio := math.Pow(10, float64(precision))
	return math.Round(val*ratio) / ratio
}

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

	p, _ := ah.ProfileServices.GetProfileByUserId(c.Get(user_id_key).(int))

	user, _ := ah.UserServices.CheckID(c.Get(user_id_key).(int))

	if c.Request().Method == "POST" {
		t := c.FormValue("t")
		if t == "pfpupdate" {
			log.Print(t)
			form, err := c.MultipartForm()
			if err != nil {
				errs["avatar"] = "Error uploading file"
			}

			files := form.File["profile_picture"]

			if len(files) == 0 {
				errs["avatar"] = "No file uploaded"
			} else {

				av := files[0]

				src, err := av.Open()
				if err != nil {
					errs["avatar"] = "Error opening file"
				}

				defer src.Close()

				s := String(4)
				filename := fmt.Sprintf("./public/avatars/%s_%s.img", c.Get(user_name_key).(string), s)
				finalurl := fmt.Sprintf("/av/%s_%s.img", c.Get(user_name_key).(string), s)

				dst, err := os.Create(filename)

				if err != nil {
					errs["avatar"] = "Error saving profile picture"
				}

				defer dst.Close()

				if _, err = io.Copy(dst, src); err != nil {
					return err
				}

				formdata["profile_picture"] = finalurl
				p.ProfilePicture = finalurl

				ah.ProfileServices.UpdateProfilePicture(user, finalurl)
			}
		} else if t == "accupdate" {
			mail := c.FormValue("email")
			username := c.FormValue("username")

			if mail != user.Email {
				if _, err := ah.UserServices.CheckEmail(mail); err == nil {
					errs["email"] = "Email already exists"
				} else if mail == "" {
					errs["email"] = "Mail Cannot be empty"
				} else {
					ah.UserServices.UpdateEmail(mail, user.ID)
					user.Email = mail
				}
			}

			if username != user.Username {
				if _, err := ah.UserServices.CheckUsername(username); err == nil {
					errs["username"] = "Username already exists"
				} else if username == "" {
					errs["username"] = "Username Cannot be empty"
				} else {
					ah.UserServices.UpdateUsername(username, user.ID)
					user.Username = username
				}
			}
		} else if t == "profileupdate" {
			log.Print("profileupdate")
			weight, err := strconv.ParseFloat(c.FormValue("weight"), 64)
			weightunit := c.FormValue("weightunit")
			if err != nil {
				errs["weight"] = "Invalid weight"
			}

			height, err := strconv.ParseFloat(c.FormValue("height"), 64)
			heightunit := c.FormValue("heightunit")
			if err != nil {
				errs["height"] = "Invalid height"
			}

			age := c.FormValue("dob")

			if !ValidateDate(age) {
				errs["dob"] = "Invalid Date"
			}

			p.Weight = weight
			p.Height = height
			p.Birthday = age

			p.HeightUnit = heightunit
			p.WeightUnit = weightunit

			if len(errs) == 0 {
				ah.ProfileServices.UpdateProfile(c.Get(user_id_key).(int), height, weight, age, heightunit, weightunit)
			}
		}
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
