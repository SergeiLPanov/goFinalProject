package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

const schema = `
CREATE TABLE scheduler (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    date CHAR(8) NOT NULL DEFAULT "",
    title VARCHAR(256) NOT NULL DEFAULT "",
    comment TEXT NOT NULL DEFAULT "",
    repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX scheduler_date ON scheduler (date);
`

var db *sql.DB

func Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	database, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("не удалось открыть БД: %w", err)
	}

	if install {
		if _, err := database.Exec(schema); err != nil {
			database.Close()
			return fmt.Errorf("не удалось создать схему БД: %w", err)
		}
	}

	db = database
	return nil
}

func Close() error {
	if db == nil {
		return nil
	}
	return db.Close()
}
