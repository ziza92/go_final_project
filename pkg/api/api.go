package api

import (
	"net/http"
)

// API хендлеры
func Init(pass string) {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler)      // Обработчик для вычисления следующей даты
	http.HandleFunc("/api/task", auth(taskHandler, pass)) // Обработчик для работы с задачами
	http.HandleFunc("/api/tasks", auth(tasksHandler, pass))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler, pass))

}
