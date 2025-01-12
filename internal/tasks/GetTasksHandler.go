package tasks

import (
	"encoding/json"
	"go_final_project/internal/db"
	"log"
	"net/http"
	"time"
)

// Обработчик для получения задач
func GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	// Установка заголовка Content-Type
	w.Header().Set("Content-Type", "application/json")

	// Получение параметра поиска
	search := r.URL.Query().Get("search")

	// Проверка, является ли search датой
	if search != "" {
		if _, err := time.Parse("02.01.2006", search); err == nil {
			// Если это дата, преобразуем в нужный формат
			date := search[6:10] + search[3:5] + search[0:2] // Преобразование 08.02.2024 в 20240208
			err = db.DB.Select(&tasks, "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50", date)
			if err != nil {
				log.Printf("Ошибка преобразования даты в нужный формат: %s\n", err)
			}
		} else {
			// Иначе выполняем поиск по заголовку и комментарию
			searchPattern := "%" + search + "%"
			err = db.DB.Select(&tasks, "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 50", searchPattern, searchPattern)
			if err != nil {
				log.Printf("Ошибка поиска по заголовку и комментарию: %s\n", err)
			}
		}
	} else {
		// Если search пустой, просто получаем все задачи
		err = db.DB.Select(&tasks, "SELECT * FROM scheduler ORDER BY date LIMIT 50")
	}

	// Обработка ошибок
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	// Если задач нет, возвращаем пустой список
	if tasks == nil {
		tasks = []Task{}
	}

	// Формируем ответ
	response := TasksResponse{Tasks: tasks}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}
