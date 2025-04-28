package main

import (
	"log"
	"os"
	"strconv"

	"golf/pkg/api"
	"golf/pkg/db"
	"golf/pkg/server"
	"golf/tests"
)

func main() {
	// определение пути к файлу (со звездочкой)
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = tests.DBFile
	}

	db, err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка создания базы данных: %e\n", err)
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		// Используем порт по умолчанию
		port = strconv.Itoa(tests.Port)
	}

	//инициализация обработчиков
	api.Init()

	err = server.StartServer(port)
	if err != nil {
		log.Fatal(err)
	}
}
