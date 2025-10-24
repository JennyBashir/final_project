package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

var schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(12),
    comment TEXT,
    repeat VARCHAR(128)
    );
CREATE INDEX scheduler_date ON scheduler (date);`

func Init(dbFile string) error {

	file := os.Getenv("TODO_DBFILE")
	if file == "" {
		file = dbFile
	}

	var install bool

	_, err := os.Stat(file)
	install = os.IsNotExist(err)

	db, err = sql.Open("sqlite", file)
	if err != nil {
		return fmt.Errorf("file %s opening error %w", file, err)
	}
	//defer db.Close()

	if install {
		_, err := db.Exec(schema)
		if err != nil {
			return fmt.Errorf("error creating the table: %w", err)
		}
	}
	return nil
}
