package server

import (
	"fmt"
	"log"
	"net/http"
)

const WebDir = "./web"

func StartServer(port string) error {

	address := fmt.Sprintf(":%s", port)

	http.Handle("/", http.FileServer(http.Dir(WebDir)))

	log.Printf("Сервер запущен на порту: %s\n", address)

	err := http.ListenAndServe(address, nil)

	return err
}
