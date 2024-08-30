package handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages/journal"
)

type EntryService interface {
	GetAllEntriesByUser(id int) []services.Entry
}

type JournalHandler struct {
	ProfileServices ProfileService
	EntryServices   EntryService
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func (jh *JournalHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	// isError = false
	jourView := journal.Journal(fromProtected)
	c.Set("ISERROR", false)

	return renderView(c, journal.JournalIndex(
		"Journal",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		jourView,
	))
}

func (jh *JournalHandler) DayHandler(c echo.Context) error {

	month, err := strconv.Atoi(c.Param("month"))
	year, err := strconv.Atoi(c.Param("year"))
	date, err := strconv.Atoi(c.Param("date"))

	formdata := make(map[string]string)

	monthname := getMonthName(month)

	isok, err := isNotFuture(year, month)

	if monthname == "0" || !isok {
		return c.Redirect(http.StatusFound, "/")
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date"})
	}

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}

	entries := jh.EntryServices.GetAllEntriesByUser(c.Get(user_id_key).(int))

	if c.Request().Method == "POST" {
		value, err := strconv.Atoi(c.FormValue("value"))
		typ := c.FormValue("type")
		stat := c.FormValue("desc")

		if err != nil {
			formdata["error"] = "Invalid values detected"
		}

		var allowed []string
		allowed = append(allowed, "hr", "bp", "sp", "height", "weight", "sugar")

		if !stringInSlice(typ, allowed) {
			formdata["error"] = "Invalid values detected"
		}

		log.Println(typ, value, stat, formdata)
	}

	jourView := journal.Day(fromProtected, entries, formdata)
	c.Set("ISERROR", false)

	return renderView(c, journal.DayIndex(
		fmt.Sprintf("%d/%d/%d", year, month, date),
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		jourView,
	))
}

func daysInMonth(year int, month int) int {
	return time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, 1, -1).Day()
}

func getMonthName(month int) string {
	if month < 1 || month > 12 {
		return "0"
	}
	return time.Month(month).String()
}

func isNotFuture(year int, month int) (bool, error) {
	inputDate := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	now := time.Now()
	currentDate := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, time.UTC)
	if inputDate.After(currentDate) {
		return false, errors.New("Date is in the future")
	}

	return true, nil
}

func getDayOfWeek(year, month, day int) int {
	date := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
	return int(date.Weekday())
}

func (jh *JournalHandler) MonthHandler(c echo.Context) error {
	month, err := strconv.Atoi(c.Param("month"))
	year, err := strconv.Atoi(c.Param("year"))

	monthname := getMonthName(month)

	isok, err := isNotFuture(year, month)

	if monthname == "0" || !isok {
		return c.Redirect(http.StatusFound, "/")
	}

	if err != nil {
		return c.Redirect(http.StatusFound, "/")
	}

	m := make([]map[string]string, daysInMonth(year, month))

	for i := 0; i < daysInMonth(year, month); i++ {
		m[i] = make(map[string]string)
		m[i]["date"] = strconv.Itoa(i + 1)
	}

	nm := month
	ny := year

	if month == 12 {
		ny += 1
		nm = 1
	} else {
		nm += 1
	}

	pm := month
	py := year

	if month == 1 {
		py -= 1
		pm = 12
	} else {
		pm -= 1
	}

	extras := getDayOfWeek(year, month, 1)
	e2 := 7 - getDayOfWeek(year, month, daysInMonth(year, month))
	if getDayOfWeek(year, month, daysInMonth(year, month)) == 0 {
		e2 = 0
	}

	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	// isError = false
	jourView := journal.Month(fromProtected, monthname, year, m, extras, e2, nm, ny, pm, py, month)
	c.Set("ISERROR", false)

	return renderView(c, journal.MonthIndex(
		"Journal",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		jourView,
	))
}
