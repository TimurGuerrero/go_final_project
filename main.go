package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	// начало Тестового кода проверок
	now, _ := time.Parse("20060102", "20240126")
	fmt.Println(NextDate(now, "16890220", "y"))

	// конец тестового кода проверок

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

	http.HandleFunc("/api/nextdate", NextDateHandler)
	http.HandleFunc("/api/task", TaskHandler)      // Обработчик для задач
	http.HandleFunc("/api/task/done", TaskHandler) // Обработчик для отметки задачи как выполненной
	http.HandleFunc("/api/tasks", GetTasksHandler)

	log.Println("Сервер запущен на порту", port)
	err = http.ListenAndServe(":"+port, nil)
	if err != nil {
		log.Fatal(err)
	}

}
