package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golf/pkg/db"
)

// Обработать запрос на обновление задачи
func updateTaskHandler(w http.ResponseWriter, r *http.Request) {

	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("failed to decode JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	if task.ID == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	if task.Title == "" {
		writeJson(w, map[string]string{"error": "title is required"}, http.StatusBadRequest)
		return
	}

	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}
	if err := db.UpdateTask(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
