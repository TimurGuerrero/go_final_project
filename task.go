package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3" // Импортируем драйвер SQLite
)

type task struct {
	ID      string `json:"id"`
	Date    string `json:"date"`
	Title   string `json:"title"`
	Comment string `json:"comment,omitempty"`
	Repeat  string `json:"repeat,omitempty"`
}

type TasksResponse struct {
	Tasks []task `json:"tasks"`
}

// jsonError формирует JSON-ответ с ошибкой
func jsonError(message string) string {
	errorResponse, _ := json.Marshal(map[string]string{"error": message})
	return string(errorResponse)
}

// Обработчик для задач
// Обработчик для задач
func TaskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		if r.URL.Path == "/api/task/done" {
			MarkTaskAsDoneHandler(w, r)
		} else {
			AddTaskHandler(w, r)
		}
	case http.MethodGet:
		GetTaskHandler(w, r) // Обработчик для получения задачи
	case http.MethodDelete:
		DeleteTaskHandler(w, r)
	case http.MethodPut:
		UpdateTaskHandler(w, r) // Обработчик для обновления задачи
	default:
		http.Error(w, jsonError("Метод не поддерживается"), http.StatusMethodNotAllowed)
	}
}

var tasks []task
var err error

// Обработчик для добавления задачи
func AddTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	var newTask task
	err := json.NewDecoder(r.Body).Decode(&newTask)
	if err != nil {
		log.Printf("Ошибка десериализации JSON: %v", err)
		http.Error(w, jsonError("Ошибка десериализации JSON"), http.StatusBadRequest)
		return
	}

	// Проверка обязательного поля title
	if newTask.Title == "" {
		http.Error(w, jsonError("Не указан заголовок задачи"), http.StatusBadRequest)
		return
	}

	// Проверка формата даты
	if newTask.Date == "" {
		newTask.Date = time.Now().Format("20060102")
	} else {
		_, err = time.Parse("20060102", newTask.Date)
		if err != nil {
			http.Error(w, jsonError("Дата представлена в неправильном формате"), http.StatusBadRequest)
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
				http.Error(w, jsonError("Неправильный формат правила повторения"), http.StatusBadRequest)
				return
			}
			newTask.Date = nextDate
		}
	}

	// Сохранение задачи в БД
	id, err := saveTaskToDB(newTask)
	if err != nil {
		http.Error(w, jsonError("Ошибка при добавлении задачи"), http.StatusInternalServerError)
		return
	}

	// Возврат идентификатора созданной задачи
	response := map[string]interface{}{"id": id}
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(response)
}

// Обработчик для получения задачи по идентификатору

func GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор"), http.StatusBadRequest)
		return
	}

	var t task
	err := DB.Get(&t, "SELECT * FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(t); err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
	}
}

// Обработчик для обновления задачи
// Обработчик для обновления задачи
// Обработчик для обновления задачи
func UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var updatedTask task
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
	}

	// Проверка существования задачи перед обновлением
	var existingTask task
	err = DB.Get(&existingTask, "SELECT * FROM scheduler WHERE id = ?", updatedTask.ID)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	// Обновление задачи в БД
	_, err = DB.Exec("UPDATE scheduler SET date = ?, title = ?, comment = ?, repeat = ? WHERE id = ?",
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
			err = DB.Select(&tasks, "SELECT * FROM scheduler WHERE date = ? ORDER BY date LIMIT 50", date)
			if err != nil {
				log.Printf("Ошибка преобразования даты в нужный формат: %s\n", err)
			}
		} else {
			// Иначе выполняем поиск по заголовку и комментарию
			searchPattern := "%" + search + "%"
			err = DB.Select(&tasks, "SELECT * FROM scheduler WHERE title LIKE ? OR comment LIKE ? ORDER BY date LIMIT 50", searchPattern, searchPattern)
			if err != nil {
				log.Printf("Ошибка поиска по заголовку и комментарию: %s\n", err)
			}
		}
	} else {
		// Если search пустой, просто получаем все задачи
		err = DB.Select(&tasks, "SELECT * FROM scheduler ORDER BY date LIMIT 50")
	}

	// Обработка ошибок
	if err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}

	// Если задач нет, возвращаем пустой список
	if tasks == nil {
		tasks = []task{}
	}

	// Формируем ответ
	response := TasksResponse{Tasks: tasks}
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, jsonError(err.Error()), http.StatusInternalServerError)
		return
	}
}

// Обработчик для отметки задачи как выполненной
func MarkTaskAsDoneHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор задачи"), http.StatusBadRequest)
		return
	}

	var t task
	err := DB.Get(&t, "SELECT * FROM scheduler WHERE id = ?", id)
	if err != nil {
		http.Error(w, jsonError("Задача не найдена"), http.StatusNotFound)
		return
	}

	if t.Repeat != "" {
		// Если задача периодическая, рассчитываем следующую дату
		currentDate, _ := time.Parse("20060102", t.Date)
		nextDate, err := NextDate(currentDate, t.Date, t.Repeat)
		if err != nil {
			http.Error(w, jsonError("Ошибка при расчете следующей даты"), http.StatusInternalServerError)
			return
		}
		t.Date = nextDate

		// Обновляем задачу в базе данных
		_, err = DB.Exec("UPDATE scheduler SET date = ? WHERE id = ?", t.Date, id)
		if err != nil {
			http.Error(w, jsonError("Ошибка при обновлении задачи"), http.StatusInternalServerError)
			return
		}
	} else {
		// Если задача одноразовая, удаляем ее
		_, err = DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
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

// Обработчик для удаления задачи
// Обработчик для удаления задачи
func DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Получен запрос: %s %s", r.Method, r.URL.Path)

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, jsonError("Не указан идентификатор задачи"), http.StatusBadRequest)
		return
	}

	// Попытка удалить задачу из базы данных
	result, err := DB.Exec("DELETE FROM scheduler WHERE id = ?", id)
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
