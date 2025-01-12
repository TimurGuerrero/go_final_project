package tasks

import (
	"encoding/json"
	"go_final_project/internal/dates"
	"go_final_project/internal/db"
	"log"
	"net/http"
	"time"
)

// Обработчик для отметки задачи как выполненной
func MarkTaskAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор задачи"), http.StatusBadRequest)
		return
	}

	var t Task
	err := db.DB.Get(&t, "SELECT * FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	if t.Repeat != "" {
		// Если задача периодическая, рассчитываем следующую дату
		currentDate, _ := time.Parse("20060102", t.Date)
		nextDate, err := dates.NextDate(currentDate, t.Date, t.Repeat)
		if err != nil {
			http.Error(w, jsonError("Ошибка при расчете следующей даты"), http.StatusInternalServerError)
			return
		}
		t.Date = nextDate

		// Обновляем задачу в базе данных
		_, err = db.DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", t.Date, id)
		if err != nil {
			http.Error(w, jsonError("Ошибка при обновлении задачи"), http.StatusInternalServerError)
			return
		}
	} else {
		// Если задача одноразовая, удаляем ее
		_, err = db.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
		if err != nil {
			http.Error(w, jsonError("Ошибка при удалении задачи"), http.StatusInternalServerError)
			return
		}
	}

	// Возвращаем пустой JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
