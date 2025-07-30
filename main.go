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

	log.Fatal(http.ListenAndServe(":8000", handlers.JsonContentTypeMiddleware(router)))
}