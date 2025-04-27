package main

import (
	"golf/pkg/api"
	"golf/pkg/db"
	"golf/pkg/server"
	"golf/tests"
	"log"
	"os"
)

func main() {
	// определение пути к файлу (со звездочкой)
	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = tests.DBFile
	}

	err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка создания базы данных: %e\n", err)
	}

	//инициализация обработчиков
	api.Init()

	err = server.StartServer()
	if err != nil {
		log.Fatal(err)
	}
}
