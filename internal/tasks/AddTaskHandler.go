package tasks

import (
	"encoding/json"
	"go_final_project/internal/dates"
	"log"
	"net/http"
	"time"
)

// Обработчик для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	var newTask Task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		http.Error(w, jsonError("Ошибка десериализации JSON"), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if newTask.Title == "" {
		http.Error(w, jsonError("Не указан заголовок задачи"), http.StatusBadRequest)
		log.Printf("Не указан заголовок задачи")
		return
	}

	// Проверка формата даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	} else {
		_, err = time.Parse("20060102", newTask.Date)
		if err != nil {
			http.Error(w, jsonError("Дата представлена в неправильном формате"), http.StatusBadRequest)
			log.Printf("Дата предоставлена в верном формате")
			return
		}
	}

	// Проверка даты
	today := time.Now().Format("20060102")
	if newTask.Date < today {
		if newTask.Repeat == "" {
			log.Printf("Формат правила повторения отсутствует")
			newTask.Date = today
		} else {
			currentDate, _ := time.Parse("20060102", today)
			nextDate, err := dates.NextDate(currentDate, newTask.Date, newTask.Repeat)
			if err != nil {
				http.Error(w, jsonError("Неправильный формат правила повторения"), http.StatusBadRequest)
				log.Printf("Формат правила повторения неверный")
				return
			}
			newTask.Date = nextDate
		}
	}

	// Сохранение задачи в БД
	id, err := SaveTaskToDB(newTask)
	if err != nil {
		http.Error(w, jsonError("Ошибка при добавлении задачи"), http.StatusInternalServerError)
		return
	}

	// Возврат идентификатора созданной задачи
	response := map[string]interface{}{"id": id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}
