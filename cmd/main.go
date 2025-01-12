package main

import (
	"log"
	"net/http"
	"os"

	"go_final_project/internal/auth"
	"go_final_project/internal/dates"
	"go_final_project/internal/db"
	"go_final_project/internal/tasks"
)

func main() {

	db.DB = db.InitDB()
	defer db.DB.Close()

	// Определяем порт из переменной окружения или используем значение по умолчанию
	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = "7540" // Порт по умолчанию
	}

	// Указываем директорию для статических файлов
	webDir := "./web"

	// Настраиваем обработчик для статических файлов
	http.Handle("/", http.FileServer(http.Dir(webDir)))

	// Запускаем сервер
	http.HandleFunc("/api/signin", auth.SigninHandler)
	http.HandleFunc("/api/nextdate", dates.NextDateHandler) // Обработчик для следующей даты
	http.HandleFunc("/api/task", tasks.TaskHandler)         // Обработчик для задач
	http.HandleFunc("/api/task/done", tasks.TaskHandler)    // Обработчик для отметки задачи как выполненной
	http.HandleFunc("/api/tasks", tasks.GetTasksHandler)    // Обработчик для получения задач

	log.Println("Сервер запущен на порту", port)
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
