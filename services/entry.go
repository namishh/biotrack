package services

import "github.com/namishh/biotrack/database"

type Entry struct {
	ID        int     `json:"id"`
	CreatedAt string  `json:"created_at"`
	Type      string  `json:"type"`
	Status    string  `json:"status"`
	CreatedBy int     `json:"created_by"`
	Value     float64 `json:"value"`
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

func (es *EntryServices) CreateEntry(user int, typ string, status string, value float64) error {

	return nil
}

func (es *EntryServices) GetAllEntriesByUser(id int) []Entry {
	stmt := `SELECT id, type, status, created_by, value, created_at from entry where created_by = ?`
	rows, err := es.EntryStore.DB.Query(stmt, id)

	var entries []Entry

	if err != nil {
		return entries
	}

	for rows.Next() {
		var e Entry
		err := rows.Scan(&e.ID, &e.Type, &e.Status, &e.CreatedBy, &e.Value, &e.CreatedAt)

		if err != nil {
			var ee []Entry
			return ee
		}

		entries = append(entries, e)
	}

	return entries
}
