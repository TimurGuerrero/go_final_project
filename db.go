package main

import (
	"log"
	"os"
	"path/filepath"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	DB     *sqlx.DB // Глобальная переменная для хранения подключения к базе данных
	DBFile string
)

// InitDB инициализирует базу данных и создает таблицу, если она отсутствует.
func InitDB() *sqlx.DB {
	// Получение текущей рабочей директории
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("Error getting current working directory: %v", err)
	}

	// Формирование пути к файлу базы данных в текущей рабочей директории
	dbFile := filepath.Join(cwd, "scheduler.db")
	DBFile = dbFile

	log.Printf("Using database file: %s", dbFile)

	// Проверка существования файла базы данных
	_, err = os.Stat(dbFile)

	var install bool
	if os.IsNotExist(err) {
		// Если файл не существует, установка флага install в true
		install = true
		log.Println("Database file does not exist. It will be created.")
		// Явное создание файла базы данных
		file, err := os.Create(dbFile)
		if err != nil {
			log.Fatalf("Error creating database file: %v", err)
		}
		file.Close() // Закрываем файл сразу после создания
	} else if err != nil {
		log.Fatalf("Error checking database file: %v", err)
	}

	// Подключение к базе данных
	log.Println("Connecting to database...")
	db, err := sqlx.Connect("sqlite3", dbFile)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	log.Println("Connected to database successfully.")

	// Установка таймаута
	if _, err := db.Exec(`PRAGMA busy_timeout = 10000;`); err != nil { // 10 секунд
		log.Fatalf("Error setting busy timeout: %v", err)
	}

	if install {
		// Если база данных не существует, создаем таблицу и индекс
		log.Println("Creating table 'scheduler'...")

		// Начинаем транзакцию
		tx := db.MustBegin()
		if _, err := tx.Exec(`CREATE TABLE scheduler (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			date CHAR(8) NOT NULL DEFAULT '',
			title VARCHAR(128) NOT NULL DEFAULT '',
			comment VARCHAR(128) NOT NULL DEFAULT '',
			repeat VARCHAR(128) NOT NULL DEFAULT ''
		)`); err != nil {
			log.Fatalf("Error creating table: %v", err)
		}
		log.Println("Table 'scheduler' created successfully.")

		// Создание индекса
		log.Println("Creating index 'idx_date'...")
		if _, err := tx.Exec(`CREATE INDEX IF NOT EXISTS idx_date ON scheduler (date)`); err != nil {
			log.Fatalf("Error creating index: %v", err)
		}

		// Завершение транзакции
		if err := tx.Commit(); err != nil {
			log.Fatalf("Error committing transaction: %v", err)
		}
		log.Println("Index 'idx_date' created successfully.")
	} else {
		log.Println("Database already exists. No need to create tables.")
	}

	// Возврат соединения с базой данных
	return db
}
