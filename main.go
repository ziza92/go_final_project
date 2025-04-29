package main

import (
	"log"
	"os"

	"golf/pkg/api"
	"golf/pkg/db"
	"golf/pkg/server"
)

const (
	PORT   = "7540"
	DBFILE = "./scheduler.db"
)

func main() {
	// определение пути к файлу (со звездочкой)
	dbFile := os.Getenv("TODO_DBFILE")
	if len(dbFile) == 0 {
		dbFile = DBFILE
	}

	db, err := db.Init(dbFile)
	if err != nil {
		log.Fatalf("Ошибка создания базы данных: %e\n", err)
	}
	defer db.Close()

	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		// Используем порт по умолчанию
		port = PORT
	}

	pass := os.Getenv("TODO_PASSWORD")
	if len(pass) == 0 {
		pass = ""
	}

	//инициализация обработчиков
	api.Init(pass)

	err = server.StartServer(port)
	if err != nil {
		log.Fatal(err)
	}
}
