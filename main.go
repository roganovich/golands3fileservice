package main

import (
	"log"
	"github.com/joho/godotenv"
	"os"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/swaggo/http-swagger"
	_ "golands3fileservice/docs"
	"golands3fileservice/pkg/database"
	"golands3fileservice/pkg/handlers"

)

// @title Golang S3 Fileservice API
// @description Простое приложение на GO, предназначенное для работы с S3 сервером.
// @version 1.0
// @host localhost:8080
// @BasePath /api
func main() {
	// Загружаем .env файл
	err := godotenv.Load(".env.local")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// InitDB
	dataSourceName := os.Getenv("DATABASE_URL")
	database.InitDB(dataSourceName)

	// Регистрация маршрутов
	router := mux.NewRouter()

	// Применяем middleware CORS ко всем роутам
	router.Use(handlers.CORS)

	// Swagger
	router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// Участники
	router.HandleFunc("/api/users", handlers.AuthAdminMiddleware(handlers.GetUsers())).Methods("GET")
	router.HandleFunc("/api/users/{id}", handlers.AuthAdminMiddleware(handlers.GetUser())).Methods("GET")

	// Кабинет
	router.HandleFunc("/api/auth/info", handlers.AuthMiddleware(handlers.InfoUser())).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/auth/create", handlers.CreateUser()).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/update", handlers.AuthMiddleware(handlers.UpdateUser())).Methods("PUT", "OPTIONS")
	router.HandleFunc("/api/auth/login", handlers.Login()).Methods("POST", "OPTIONS")
	router.HandleFunc("/api/auth/refresh", handlers.AuthMiddleware(handlers.Refresh())).Methods("POST", "OPTIONS")

	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}