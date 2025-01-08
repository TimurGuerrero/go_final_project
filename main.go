package main

import (
	"log"
	"net/http"
	"os"
)

func main() {

	DB = InitDB()
	defer DB.Close()

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
	http.HandleFunc("/api/signin", SignInHandler)
	http.HandleFunc("/api/nextdate", NextDateHandler)    // Обработчик для следующей даты
	http.HandleFunc("/api/task", auth(TaskHandler))      // Обработчик для задач
	http.HandleFunc("/api/task/done", auth(TaskHandler)) // Обработчик для отметки задачи как выполненной
	http.HandleFunc("/api/tasks", auth(GetTasksHandler)) // Обработчик для получения задач

	log.Println("Сервер запущен на порту", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
