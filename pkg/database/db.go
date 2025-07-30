package database

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

// DB - глобальная переменная для соединения с базой данных
var DB *sql.DB

// InitDB - функция для инициализации соединения с базой данных
func InitDB(dataSourceName string) {
	var dbError error
	DB, dbError = sql.Open("postgres", dataSourceName)
	if dbError != nil {
		log.Fatalf("Не удалось открыть соединение с базой данных: %v", dbError)
	}

	if dbError = DB.Ping(); dbError != nil {
		log.Fatalf("Не удалось установить соединение с базой данных: %v", dbError)
	}

	log.Println("Соединение с базой данных установлено.")
}