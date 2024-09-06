package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/services"
	"github.com/namishh/biotrack/views/pages/journal"
)

type EntryService interface {
	GetAllEntriesByUser(id int) ([]services.Entry, error)
	GetAllEntriesByDate(id int, month, day, year int) ([]services.Entry, error)
	GetAllEntriesByMonth(id int, month, year int) ([]services.Entry, error)
	CreateEntry(user int, typ string, status string, value float64, month, day, year int) error
	DeleteEntry(id int) error
	GetFormattedEntriesByUser(userID int) (string, error)
	GetEntryByID(id int) (services.Entry, error)
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

func groupEntriesByType(entries []services.Entry) map[string][]services.Entry {
	result := make(map[string][]services.Entry)

	// First, group all entries by type
	for _, entry := range entries {
		result[entry.Type] = append(result[entry.Type], entry)
	}

	// Then, for each type, sort the entries and keep only the last 10 (or all if less than 10)
	for entryType, typeEntries := range result {
		// Sort entries by date in increasing order
		sort.Slice(typeEntries, func(i, j int) bool {
			dateI := time.Date(typeEntries[i].Year, time.Month(typeEntries[i].Month), typeEntries[i].Day, 0, 0, 0, 0, time.UTC)
			dateJ := time.Date(typeEntries[j].Year, time.Month(typeEntries[j].Month), typeEntries[j].Day, 0, 0, 0, 0, time.UTC)
			return dateI.Before(dateJ)
		})

		// Keep only the last 10 entries (or all if less than 10)
		if len(typeEntries) > 10 {
			result[entryType] = typeEntries[len(typeEntries)-10:]
		} else {
			result[entryType] = typeEntries
		}
	}

	return result
}

func (jh *JournalHandler) HomeHandler(c echo.Context) error {
	fromProtected, ok := c.Get("FROMPROTECTED").(bool)
	if !ok {
		return errors.New("invalid type for key 'FROMPROTECTED'")
	}
	entries, err := jh.EntryServices.GetAllEntriesByUser(c.Get(user_id_key).(int))
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get entries"})
	}

	groupedEntries := groupEntriesByType(entries)
	profile, err := jh.ProfileServices.GetProfileByUserId(c.Get(user_id_key).(int))

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Failed to get entries"})
	}
	// isError = false
	jourView := journal.Journal(fromProtected, groupedEntries, profile)
	c.Set("ISERROR", false)

	return renderView(c, journal.JournalIndex(
		"Journal",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		jourView,
	))
}

func isCurrentDay(year, month, date int) bool {
	// Create a time.Time object from the given date
	givenDate := time.Date(year, time.Month(month), date, 0, 0, 0, 0, time.Local)

	// Get the current date
	currentDate := time.Now()

	// Compare year, month, and day
	return givenDate.Year() == currentDate.Year() &&
		givenDate.Month() == currentDate.Month() &&
		givenDate.Day() == currentDate.Day()
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

	entries, err := jh.EntryServices.GetAllEntriesByDate(c.Get(user_id_key).(int), month, date, year)

	if err != nil {
		formdata["error"] = "Error Fetching Entries"
	}

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

		// Validate data
		if value <= 0 {
			formdata["error"] = "Invalid values detected"
		}

		if typ == "sp" && value > 100 {
			formdata["value"] = "Invalid values Detected"
		}

		if typ == "hr" && value > 200 {
			formdata["value"] = "Invalid values Detected"
		}

		// Update profile's height and weight if the entry is of type height or weight

		if typ == "height" && isCurrentDay(year, month, date) {
			err = jh.ProfileServices.UpdateProfileHeight(c.Get(user_id_key).(int), value)
		}

		if typ == "weight" && isCurrentDay(year, month, date) {
			err = jh.ProfileServices.UpdateProfileWeight(c.Get(user_id_key).(int), value)
		}

		if !stringInSlice(typ, allowed) {
			formdata["error"] = "Error Creating Entry"
		}

		if len(formdata) < 1 {
			err = jh.EntryServices.CreateEntry(c.Get(user_id_key).(int), typ, stat, float64(value), month, date, year)
			entries, err = jh.EntryServices.GetAllEntriesByDate(c.Get(user_id_key).(int), month, date, year)
		}

		if err != nil {
			formdata["error"] = "Error Creating Entry"
		}
	}

	jourView := journal.Day(fromProtected, entries, formdata, year, month, date)
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

	// make a list of list of Entries in order of the days of the month

	en := make([][]services.Entry, daysInMonth(year, month))

	for i := 0; i < daysInMonth(year, month); i++ {
		en[i] = make([]services.Entry, 0)
		entries, err := jh.EntryServices.GetAllEntriesByDate(c.Get(user_id_key).(int), month, i+1, year)
		if err != nil {
			return c.Redirect(http.StatusFound, "/")
		}

		for _, entry := range entries {
			en[i] = append(en[i], entry)
		}
	}

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
	jourView := journal.Month(fromProtected, monthname, year, m, extras, e2, nm, ny, pm, py, month, en)
	c.Set("ISERROR", false)

	return renderView(c, journal.MonthIndex(
		"Journal",
		"",
		fromProtected,
		c.Get("ISERROR").(bool),
		jourView,
	))
}

func (jh *JournalHandler) DeleteHandler(c echo.Context) error {
	month, err := strconv.Atoi(c.Param("month"))
	year, err := strconv.Atoi(c.Param("year"))
	date, err := strconv.Atoi(c.Param("date"))
	id, err := strconv.Atoi(c.Param("id"))

	monthname := getMonthName(month)

	isok, err := isNotFuture(year, month)

	if monthname == "0" || !isok {
		return c.Redirect(http.StatusFound, "/")
	}

	en, err := jh.EntryServices.GetEntryByID(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delete id"})
	}

	if en.CreatedBy != c.Get(user_id_key).(int) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delete id"})
	}

	if en.Day != date || en.Month != month || en.Year != year {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid delete id"})
	}

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid date"})
	}

	err = jh.EntryServices.DeleteEntry(id)

	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Error Deleting"})
	}

	return c.Redirect(http.StatusSeeOther, fmt.Sprintf("/journal/%d/%d/%d", year, month, date))
}

func (jh *JournalHandler) NewHandler(c echo.Context) error {
	currentDate := time.Now()
	string := fmt.Sprintf("/journal/%d/%d/%d", currentDate.Year(), int(currentDate.Month()), currentDate.Day())

	return c.Redirect(http.StatusFound, string)
}

func (jh *JournalHandler) CalendarHandler(c echo.Context) error {
	currentDate := time.Now()
	string := fmt.Sprintf("/journal/%d/%d", currentDate.Year(), int(currentDate.Month()))

	return c.Redirect(http.StatusFound, string)
}
