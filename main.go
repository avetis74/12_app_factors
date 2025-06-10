package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/avetis74/12_app_factors/handlers"
	"github.com/avetis74/12_app_factors/storage"

	"github.com/labstack/echo/v4"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL. Пустой импорт нужен для регистрации драйвера.
)

func main() {
	// Фактор III: Конфигурация из переменных окружения
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	// Подключаемся к базе данных
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("could not connect to database: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("database is not reachable: %v", err)
	}

	log.Println("Database connection successful")

	// Создаем экземпляры наших зависимостей
	userStore := storage.NewPostgresStore(db)
	userHandler := handlers.NewUserHandler(userStore)

	e := echo.New()

	// Routes
	e.GET("/users", userHandler.GetUsers)
	e.POST("/users", userHandler.CreateUser)
	// Другие роуты...
	// e.GET("/users/:id", userHandler.GetUser)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	addr := fmt.Sprintf(":%s", port)
	e.Logger.Fatal(e.Start(addr))
}
