package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/avetis74/12_app_factors/handlers"
	"github.com/avetis74/12_app_factors/storage"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	_ "github.com/lib/pq" // Драйвер для PostgreSQL
)

func main() {
	// Фактор XI: Логи как потоки событий - настройка логирования
	log.SetOutput(os.Stdout)

	// Фактор III: Конфигурация из переменных окружения
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is not set")
	}

	redisURL := os.Getenv("REDIS_URL")
	if redisURL == "" {
		log.Fatal("REDIS_URL environment variable is not set")
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

	// Подключаемся к Redis
	redisCache, err := storage.NewRedisCache(redisURL)
	if err != nil {
		log.Fatalf("could not connect to Redis: %v", err)
	}
	defer redisCache.Close()

	// Создаем экземпляры наших зависимостей
	userStore := storage.NewPostgresStore(db)
	
	// Оборачиваем store в кеширующий слой
	cachedUserStore := storage.NewCachedUserStore(userStore, redisCache)
	
	userHandler := handlers.NewUserHandler(cachedUserStore)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/users", userHandler.GetUsers)
	e.POST("/users", userHandler.CreateUser)
	e.GET("/users/:id", userHandler.GetUser)
	e.PUT("/users/:id", userHandler.UpdateUser)
	e.DELETE("/users/:id", userHandler.DeleteUser)

	// Health check endpoint
	e.GET("/health", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "healthy"})
	})

	// Cache stats endpoint
	e.GET("/cache/stats", func(c echo.Context) error {
		// Простая проверка Redis
		ctx := context.Background()
		err := redisCache.Set(ctx, "test", "ping", time.Second)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"redis": "error",
				"error": err.Error(),
			})
		}
		
		var result string
		err = redisCache.Get(ctx, "test", &result)
		if err != nil {
			return c.JSON(http.StatusInternalServerError, map[string]string{
				"redis": "error",
				"error": err.Error(),
			})
		}

		return c.JSON(http.StatusOK, map[string]string{
			"redis": "healthy",
			"test":  result,
		})
	})

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	
	// Фактор IX: Disposability - Graceful shutdown
	go func() {
		addr := fmt.Sprintf(":%s", port)
		log.Printf("Starting server on port %s", port)
		if err := e.Start(addr); err != nil && err != http.ErrServerClosed {
			log.Fatalf("shutting down the server: %v", err)
		}
	}()

	// Ожидание сигнала завершения
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// Graceful shutdown с таймаутом
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	
	if err := e.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}
	
	log.Println("Server exited")
}
