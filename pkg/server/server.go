package server

import (
	"fmt"
	"golf/tests"
	"log"
	"net/http"
	"os"
	"strconv"
)

const WebDir = "./web"

func StartServer() error {
	port := os.Getenv("TODO_PORT")
	if len(port) == 0 {
		// Используем порт по умолчанию
		port = strconv.Itoa(tests.Port)
	}
	address := fmt.Sprintf(":%s", port)

	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	log.Printf("Сервер запущен на порту: %s\n", address)

	err := http.ListenAndServe(address, nil)

	return err
}
