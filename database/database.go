package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

type DatabaseStore struct {
	DB *sql.DB
}

func GetConnection(dbName string) (*sql.DB, error) {
	var (
		err error
		db  *sql.DB
	)

	if db != nil {
		return db, nil
	}

	db, err = sql.Open("sqlite3", dbName)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect to the database: %s", err)
	}

	log.Println("Connected Successfully to the Database")

	return db, nil
}

func CreateMigrations(DBName string, DB *sql.DB) error {
	stmt := `CREATE TABLE IF NOT EXISTS users (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(20) NOT NULL,
		password VARCHAR(20) NOT NULL,
		username VARCHAR(20) NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	_, err := DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS profile (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		level INTEGER DEFAULT 0,
		weight_unit VARCHAR(5) DEFAULT 'kg',
		height_unit VARCHAR(5) DEFAULT 'cm',
		profile_picture VARCHAR(255),
		weight FLOAT DEFAULT 0,
		height FLOAT DEFAULT 0,
		birthday DATE DEFAULT '1900-01-01',
		bio TEXT DEFAULT '',
		profile_of INT,
		FOREIGN KEY(profile_of) REFERENCES users(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS entry (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME default CURRENT_TIMESTAMP,
		type TEXT NOT NULL,
		status TEXT NOT NULL,
		created_by INT NOT NULL,
		value FLOAT NOT NULL,
		month INT NOT NULL,
		year INT NOT NULL,
		day INT NOT NULL,
		FOREIGN KEY(created_by) REFERENCES users(id)
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS charts (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			created_at DATETIME default CURRENT_TIMESTAMP,
			labels TEXT NOT NULL,
			data TEXT NOT NULL,
			chat_id INT NOT NULL,
			FOREIGN KEY(chat_id) REFERENCES chat(id)
		);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	stmt = `CREATE TABLE IF NOT EXISTS chat (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		created_at DATETIME default CURRENT_TIMESTAMP,
		sender TEXT,
		message TEXT NOT NULL
	);`

	_, err = DB.Exec(stmt)
	if err != nil {
		return fmt.Errorf("Failed to create table: %s", err)
	}

	// stmt = `delete from chat;`

	// _, err = DB.Exec(stmt)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create table: %s", err)
	// }

	return nil
}

func NewDatabaseStore(path string) (DatabaseStore, error) {
	DB, err := GetConnection(path)
	if err != nil {
		return DatabaseStore{}, err
	}

	if err := CreateMigrations(path, DB); err != nil {
		return DatabaseStore{}, err
	}

	return DatabaseStore{DB: DB}, nil
}
