package services

import (
	"encoding/json"
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

type FormattedEntry struct {
	ID        int     `json:"id"`
	CreatedAt string  `json:"created_at"`
	Status    string  `json:"status"`
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
	return err
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

func (es *EntryServices) GetAllEntriesByMonth(id int, month, year int) ([]Entry, error) {
	stmt := `SELECT id, type, status, created_by, value, created_at, month, day, year FROM entry WHERE created_by = ? AND month = ? AND year = ?`
	rows, err := es.EntryStore.DB.Query(stmt, id, month, year)
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

func (es *EntryServices) DeleteEntry(id int) error {
	query := "DELETE FROM entry WHERE id = ?"

	// Execute the delete statement
	_, err := es.EntryStore.DB.Exec(query, id)
	if err != nil {
		log.Printf("Error deleting hint with ID %d: %v", id, err)
		return err
	}
	return nil
}

func (es *EntryServices) GetEntryByID(id int) (Entry, error) {
	query := `SELECT id, type, status, created_by, value, created_at, month, day, year FROM entry WHERE id = ?`

	stmt, err := es.EntryStore.DB.Prepare(query)
	if err != nil {
		log.Print(err)
		return Entry{}, err
	}

	defer stmt.Close()
	e := Entry{}

	err = stmt.QueryRow(id).Scan(&e.ID, &e.Type, &e.Status, &e.CreatedBy, &e.Value, &e.CreatedAt, &e.Month, &e.Day, &e.Year)
	if err != nil {
		log.Print(err)
		return Entry{}, err
	}

	return e, nil
}

func (es *EntryServices) GetFormattedEntriesByUser(userID int) (string, error) {
	entries, err := es.GetAllEntriesByUser(userID)
	if err != nil {
		log.Printf("Error getting entries for user %d: %v", userID, err)
		return "", err
	}

	formattedEntries := make(map[string][]FormattedEntry)

	for _, entry := range entries {
		formattedEntry := FormattedEntry{
			ID:        entry.ID,
			CreatedAt: entry.CreatedAt,
			Status:    entry.Status,
			Value:     entry.Value,
			Month:     entry.Month,
			Year:      entry.Year,
			Day:       entry.Day,
		}

		// Append the entry to the appropriate type array
		formattedEntries[entry.Type] = append(formattedEntries[entry.Type], formattedEntry)
	}

	// Convert the formatted entries to JSON
	jsonData, err := json.MarshalIndent(formattedEntries, "", "  ")
	if err != nil {
		log.Printf("Error marshaling entries to JSON: %v", err)
		return "", err
	}

	return string(jsonData), nil
}
