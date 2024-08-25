package services

import "github.com/namishh/biotrack/database"

type Entry struct {
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
