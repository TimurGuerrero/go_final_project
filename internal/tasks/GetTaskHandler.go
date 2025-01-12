package tasks

import (
	"encoding/json"
	"go_final_project/internal/db"
	"net/http"
)

// Обработчик для получения задачи по идентификатору
func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	var t Task
	err := db.DB.Get(&t, "SELECT * FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
	}
}
