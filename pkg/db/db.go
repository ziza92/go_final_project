package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var db *sql.DB

const schema = `
			CREATE TABLE IF NOT EXISTS scheduler (
				id INTEGER PRIMARY KEY AUTOINCREMENT,
				date CHAR(8) NOT NULL DEFAULT "",
				title VARCHAR(255) NOT NULL,
				comment TEXT,
				repeat VARCHAR(128)
			);

			CREATE INDEX IF NOT EXISTS idx_date ON scheduler(date);`

// Создать таблицу, если её нет
func Init(dbFile string) (*sql.DB, error) {
	// Проверка, существует ли база данных
	_, err := os.Stat(dbFile)
	if err != nil {
		// Если база данных не существует, выводим сообщение
		fmt.Println("База данных не найдена, создаем...")
	}

	db, err = sql.Open("sqlite", dbFile)
	if err != nil {
		return nil, fmt.Errorf("невозможно открыть таблицу: %v", err)
	}

	// Создать таблицу, если её нет
	_, err = db.Exec(schema)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения schema: %v", err)
	}

	return db, nil
}
