package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func openDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "sqlite3.db")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	DB = db
	return db, nil
}

func databaseFileExists() bool {
	_, err := os.Stat("sqlite3.db")
	return err == nil
}

func setupDB() error {
	_, err := DB.Exec(`
	create table if not exists tasks 
	(id integer not null primary key, 
		title text, 
		completed boolean default false, 
		position integer)`,
	)
	if err != nil {
		return err
	}

	return nil
}

func (app *application) connectToDB() (*sql.DB, error) {
	if !databaseFileExists() {
		log.Println("Database file does not exist. Creating...")
		file, err := os.Create("sqlite3.db")
		if err != nil {
			return nil, err
		}
		file.Close()
	}

	connection, err := openDB()
	if err != nil {
		return nil, err
	}

	err = setupDB()
	if err != nil {
		return nil, err
	}

	log.Println("Connect to SQLite!")
	return connection, nil
}
