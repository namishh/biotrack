package services

import (
	"log"

	"github.com/namishh/biotrack/database"
)

type Entry struct {
	ID        int     `json:"id"`
	CreatedAt string  `json:"created_at"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	CreatedBy int     `json:"created_by"`
	Value     float64 `json:"value"`
	Month     int     `json:"month"`
	Year      int     `json:"year"`
	Day       int     `json:"day"`
}

type EntryServices struct {
	Entry      Entry
	EntryStore database.DatabaseStore
}

func NewEntryService(entry Entry, store database.DatabaseStore) *EntryServices {
	return &EntryServices{
		Entry:      entry,
		EntryStore: store,
	}
}

func (es *EntryServices) CreateEntry(user int, typ string, status string, value float64, month, day, year int) error {
	stmt := `INSERT INTO entry (value, status, type, created_by, month, day, year) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err := es.EntryStore.DB.Exec(stmt, value, status, typ, user, month, day, year)

	if err != nil {
		return err
	}

	return nil
}

func (es *EntryServices) GetAllEntriesByUser(id int) ([]Entry, error) {
	stmt := `SELECT id, type, status, created_by, value, created_at, month, day, year FROM entry WHERE created_by = ?`
    rows, err := es.EntryStore.DB.Query(stmt, id)
    if err != nil {
        log.Printf("Query failed: %v", err)
        return nil, err
    }
    defer rows.Close()

    var entries []Entry
    for rows.Next() {
        var e Entry
        err := rows.Scan(&e.ID, &e.Type, &e.Status, &e.CreatedBy, &e.Value, &e.CreatedAt, &e.Month, &e.Day, &e.Year)
        if err != nil {
            log.Printf("Scan failed: %v", err)
            return nil, err
        }
        entries = append(entries, e)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Row iteration failed: %v", err)
        return nil, err
    }

    return entries, nil
}

func (es *EntryServices) GetAllEntriesByDate(id int, month, day, year int) ([]Entry, error) {
    stmt := `SELECT id, type, status, created_by, value, created_at, month, day, year FROM entry WHERE created_by = ? AND month = ? AND day = ? AND year = ?`
    rows, err := es.EntryStore.DB.Query(stmt, id, month, day, year)
    if err != nil {
        log.Printf("Query failed: %v", err)
        return nil, err
    }
    defer rows.Close()

    var entries []Entry
    for rows.Next() {
        var e Entry
        err := rows.Scan(&e.ID, &e.Type, &e.Status, &e.CreatedBy, &e.Value, &e.CreatedAt, &e.Month, &e.Day, &e.Year)
        if err != nil {
            log.Printf("Scan failed: %v", err)
            return nil, err
        }
        entries = append(entries, e)
    }

    if err = rows.Err(); err != nil {
        log.Printf("Row iteration failed: %v", err)
        return nil, err
    }

    return entries, nil
}
