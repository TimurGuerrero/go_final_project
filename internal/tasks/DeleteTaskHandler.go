package tasks

import (
	"encoding/json"
	"go_final_project/internal/db"
	"log"
	"net/http"
)

// Обработчик для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор задачи"), http.StatusBadRequest)
		return
	}

	// Попытка удалить задачу из базы данных
	result, err := db.DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, jsonError("Ошибка при удалении задачи"), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil || rowsAffected == 0 {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	// Возвращаем пустой JSON
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{})
}
