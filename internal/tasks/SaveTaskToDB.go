package tasks

import (
	"fmt"
	"go_final_project/internal/db"
	"log"
)

// Функция для сохранения задачи в БД
func SaveTaskToDB(task Task) (int64, error) {
	if db.DB == nil {
		return 0, fmt.Errorf("database connection is nil")
	}
	res, err := db.DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
	if err != nil {
		log.Printf("Ошибка выполнения запроса: %v", err)
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		log.Printf("Ошибка получения ID: %s\n", err)
	}
	return id, nil
}
