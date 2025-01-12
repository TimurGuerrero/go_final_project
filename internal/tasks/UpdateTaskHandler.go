package tasks

import (
	"encoding/json"
	"go_final_project/internal/dates"
	"go_final_project/internal/db"
	"log"
	"net/http"
	"time"
)

// Обработчик для обновления задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTask Task
	err := json.NewDecoder(r.Body).Decode(&updatedTask)
	if err != nil {
		http.Error(w, jsonError("Ошибка десериализации JSON"), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля id
	if updatedTask.ID == "" {
		http.Error(w, jsonError("Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if updatedTask.Title == "" {
		http.Error(w, jsonError("Не указан заголовок задачи"), http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if updatedTask.Date != "" {
		_, err = time.Parse("20060102", updatedTask.Date)
		if err != nil {
			http.Error(w, jsonError("Дата представлена в неправильном формате"), http.StatusBadRequest)
			return
		}
		currentDate, _ := time.Parse("20060102", updatedTask.Date)
		_, err := dates.NextDate(currentDate, updatedTask.Date, updatedTask.Repeat)
		if err != nil {
			http.Error(w, jsonError("Неправильный формат правила повторения"), http.StatusBadRequest)
			log.Printf("Формат правила повторения неверный")
			return
		}
	}

	// Проверка существования задачи перед обновлением
	var existingTask Task
	err = db.DB.Get(&existingTask, "SELECT * FROM scheduler WHERE id = ?", updatedTask.ID)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	// Обновление задачи в БД
	_, err = db.DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
		updatedTask.Date, updatedTask.Title, updatedTask.Comment, updatedTask.Repeat, updatedTask.ID)
	if err != nil {
		http.Error(w, jsonError("Ошибка при обновлении задачи"), http.StatusInternalServerError)
		return
	}

	// Возврат пустого JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
