package api

import (
	"fmt"
	"net/http"
	"time"

	"golf/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

// Проверить формат времени
func isValidDateFormat(date string) bool {
	_, err := time.Parse("02.01.2006", date)
	return err == nil
}

// Обработчик GET
func tasksHandler(w http.ResponseWriter, r *http.Request) {
	// Получить search
	search := r.URL.Query().Get("search")

	// Распарсить дату, если она есть в запросе
	var dateSearch string
	if isValidDateFormat(search) {
		dateSearch = formatDateForSearch(search)
		search = ""
	}

	tasks, err := db.Tasks(50, search, dateSearch)
	if err != nil {
		writeJson(w, map[string]string{"Ошибка": "не удалось полученить задачи: " + err.Error()}, http.StatusInternalServerError)
		return
	}

	if tasks == nil {
		tasks = []*db.Task{}
	}

	result := make([]map[string]string, 0, len(tasks))
	for _, t := range tasks {
		result = append(result, map[string]string{
			"id":      fmt.Sprint(t.ID),
			"date":    t.Date,
			"title":   t.Title,
			"comment": t.Comment,
			"repeat":  t.Repeat,
		})
	}
	writeJson(w, map[string]any{
		"tasks": result,
	}, http.StatusOK)
}

// Преобразовать дату из dd.mm.yyyy в yyyyMMdd
func formatDateForSearch(date string) string {
	parsedDate, err := time.Parse("02.01.2006", date)
	if err != nil {
		return ""
	}
	return parsedDate.Format(ConstDate)
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, task, http.StatusOK)
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {

	id := r.FormValue("id")
	if id == "" {
		writeJson(w, map[string]string{"error": "id is required"}, http.StatusBadRequest)
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	writeJson(w, map[string]interface{}{}, http.StatusOK)
}
