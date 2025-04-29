package api

import (
	"net/http"
)

// Обработчик
func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		// Добавить задачу
		addTaskHandler(w, r)
	case http.MethodGet:
		// Получить задачу
		getTaskHandler(w, r)
	case http.MethodPut:
		// Обновить задачу
		updateTaskHandler(w, r)
	case http.MethodDelete:
		// Удалить задачу
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}
}
