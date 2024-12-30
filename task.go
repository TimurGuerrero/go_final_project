package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	// Импортируем sqlx
	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
)

type task struct {
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

// Обработчик для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Received request: %s %s", r.Method, r.URL.Path)
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	var newTask task
	err := json.NewDecoder(r.Body).Decode(&newTask)

	if err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		http.Error(w, `{"error": "Ошибка десериализации JSON"}`, http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if newTask.Title == "" {
		http.Error(w, `{"error": "Не указан заголовок задачи"}`, http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	} else {
		_, err = time.Parse("20060102", newTask.Date)
		if err != nil {
			http.Error(w, `{"error": "Дата представлена в неправильном формате"}`, http.StatusBadRequest)
			return
		}
	}

	// Проверка даты
	today := time.Now().Format("20060102")
	if newTask.Date < today {
		if newTask.Repeat == "" {
			newTask.Date = today
		} else {
			currentDate, _ := time.Parse("20060102", today)
			nextDate, err := NextDate(currentDate, newTask.Date, newTask.Repeat)
			if err != nil {
				http.Error(w, `{"error": "Неправильный формат правила повторения"}`, http.StatusBadRequest)
				return
			}
			newTask.Date = nextDate
		}
	}

	// Сохранение задачи в БД
	id, err := saveTaskToDB(newTask)
	if err != nil {
		http.Error(w, `{"error": "Ошибка при добавлении задачи"}`, http.StatusInternalServerError)
		return
	}

	// Возврат идентификатора созданной задачи
	response := map[string]interface{}{"id": id}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// Функция для сохранения задачи в БД
func saveTaskToDB(task task) (int64, error) {
	if DB == nil {
		return 0, fmt.Errorf("database connection is nil")
	}
	res, err := DB.Exec("INSERT INTO scheduler (date, title, comment, repeat) VALUES (?, ?, ?, ?)", task.Date, task.Title, task.Comment, task.Repeat)
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
