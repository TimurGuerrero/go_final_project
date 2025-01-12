package tasks

import (
	"encoding/json"
	"net/http"
)

var tasks []Task
var err error

type Task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksResponse struct {
	Tasks []Task `json:"tasks"`
}

// jsonError формирует JSON-ответ с ошибкой
func jsonError(message string) string {
	errorResponse, _ := json.Marshal(map[string]string{"error": message})
	return string(errorResponse)
}

// Обработчик для задач
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/api/task/done" {
			MarkTaskAsDoneHandler(w, r) // Обработчик для отметки задачи как выполненной
		} else {
			AddTaskHandler(w, r) // Обработчик для добавления задачи
		}
	case http.MethodGet:
		GetTaskHandler(w, r) // Обработчик для получения задачи
	case http.MethodDelete:
		DeleteTaskHandler(w, r) // Обработчик для удаления задачи
	case http.MethodPut:
		UpdateTaskHandler(w, r) // Обработчик для обновления задачи
	default:
		http.Error(w, jsonError("Метод не поддерживается"), http.StatusMethodNotAllowed)
	}
}
