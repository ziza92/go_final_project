package api

import (
	"net/http"
)

// API хендлеры
func Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDayHandler) // Обработчик для вычисления следующей даты
	http.HandleFunc("/api/task", auth(taskHandler))  // Обработчик для работы с задачами
	http.HandleFunc("/api/tasks", auth(tasksHandler))
	http.HandleFunc("/api/task/done", auth(taskDoneHandler))

}
