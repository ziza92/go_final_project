package db

import (
	"database/sql"
	"fmt"
)

// func Init(database *sql.DB) {
// 	db = database
// }

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment"`
	Repeat  string `json:"repeat"`
}

// Добавить задачу
func AddTask(task *Task) (int64, error) {
	var id int64
	query := `INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat)
	if err == nil {
		id, err = res.LastInsertId()
	}

	return id, err
}

// Обновление задачи
func UpdateTask(task *Task) error {
	query := `UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?`
	res, err := db.Exec(query, task.Date, task.Title, task.Comment, task.Repeat, task.ID)
	if err != nil {
		return err
	}

	count, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}

// Получить задачу
func GetTask(id string) (*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler WHERE id = ?`

	task := &Task{}

	err := db.QueryRow(query, id).Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
	if err != nil {
		if err == sql.ErrNoRows {
			// Вернуть ошибку, если задача не найдена
			return nil, fmt.Errorf("task not found")
		}
		// Вернуть ошибку, если строка пуста
		return nil, err
	}

	return task, nil
}

// Поиск задач
func Tasks(limit int, search string, dateSearch string) ([]*Task, error) {
	query := `SELECT id, date, title, comment, repeat FROM scheduler`

	// Если есть параметр search, добавляем условие LIKE
	if search != "" {
		query += ` WHERE title LIKE ? OR comment LIKE ?`
	}

	// Если есть параметр dateSearch (поиск по дате), добавляем условие для даты
	if dateSearch != "" {
		if search != "" {
			query += ` AND date = ?`
		} else {
			query += ` WHERE date = ?`
		}
	}

	// Сортировка по дате
	query += ` ORDER BY date LIMIT ?`

	// Параметры запроса
	var params []interface{}
	if search != "" {
		params = append(params, "%"+search+"%", "%"+search+"%")
	}
	if dateSearch != "" {
		params = append(params, dateSearch)
	}
	params = append(params, limit)

	// Выполнить запрос
	rows, err := db.Query(query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tasks []*Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Date, &task.Title, &task.Comment, &task.Repeat)
		if err != nil {
			return nil, err
		}
		tasks = append(tasks, &task)
	}

	// Вернуть пустой массив, если нет задач
	if tasks == nil {
		tasks = []*Task{}
	}

	return tasks, nil
}

func DeleteTask(id string) error {
	query := `DELETE FROM scheduler WHERE id = ?`
	res, err := db.Exec(query, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("task not found")
	}
	return nil
}

func UpdateDate(next, id string) error {
	query := `UPDATE scheduler SET date = ? WHERE id = ?`
	res, err := db.Exec(query, next, id)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count == 0 {
		return fmt.Errorf("задача не найдена")
	}
	return nil
}
