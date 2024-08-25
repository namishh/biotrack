package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/namishh/biotrack/views/pages/journal"
)

type EntryService interface {
}

type JournalHandler struct {
	ProfileServices ProfileService
	EntryServices   EntryService
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
