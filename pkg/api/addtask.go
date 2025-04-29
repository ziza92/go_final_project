package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"golf/pkg/db"
)

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Чтение тела запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка чтения тела запроса: %v", err)}, http.StatusBadRequest)
		return
	}
	log.Printf("Полученные данные: %s", body)

	// Декодирование JSON
	decoder := json.NewDecoder(bytes.NewReader(body))
	if err := decoder.Decode(&task); err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка декодирования JSON: %v", err)}, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля
	if task.Title == "" {
		writeJson(w, map[string]string{"error": "Поле 'title' обязательно"}, http.StatusBadRequest)
		return
	}

	// Проверка даты
	if err := checkDate(&task); err != nil {
		writeJson(w, map[string]string{"error": err.Error()}, http.StatusBadRequest)
		return
	}

	// Добавление задачи в базу данных
	id, err := db.AddTask(&task)
	if err != nil {
		writeJson(w, map[string]string{"error": fmt.Sprintf("Ошибка добавления задачи в базу данных: %v", err)}, http.StatusInternalServerError)
		return
	}

	// Ответ с id задачи
	writeJson(w, map[string]interface{}{"id": id}, http.StatusOK)
}

// Функция проверки и корректировки даты
func checkDate(task *db.Task) error {
	now := time.Now()

	if task.Date == "" {
		task.Date = now.Format(ConstDate)
	}

	t, err := time.Parse(ConstDate, task.Date)
	if err != nil {
		return fmt.Errorf("неподдерживаемый формат даты: %v", err)
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return fmt.Errorf("неподдерживаемый формат повторения: %v", err)
		}

		if afterNow(now, t) {
			task.Date = next
		}
	} else if afterNow(now, t) {

		task.Date = now.Format(ConstDate)
	}

	return nil
}

// утилита для json ответов
func writeJson(w http.ResponseWriter, data any, statusCode int) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(statusCode)

	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		log.Printf("ошибка сериализации JSON: %v", err)
		http.Error(w, "ошибка сериализации ответа", http.StatusInternalServerError)
	}
}
