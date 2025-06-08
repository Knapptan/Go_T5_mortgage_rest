// Package main запускает HTTP-сервер Mortgage Calculator
package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Knapptan/Go_T5_mortgage_rest/config"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/cache"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/calculator"
	myhttp "github.com/Knapptan/Go_T5_mortgage_rest/internal/http"
	"github.com/gin-gonic/gin"
)

const configFile = "config.yml"

func main() {
	// Загружаем конфигурацию из файла
	cfg, err := config.Load(configFile)
	if err != nil {
		log.Fatalf("Failed to load config from %q: %v", configFile, err)
	}

	// Инициализируем кэш, сервис калькулятора и HTTP-хендлер
	mortgageCache := cache.NewCache()
	calcService := calculator.NewService()
	handler := myhttp.NewHandler(calcService, mortgageCache)

	// Настраиваем Gin в режиме релиза и создаем роутер
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	myhttp.SetupRoutes(router, handler)

	// Создаем HTTP-сервер с таймаутами для чтения, записи и простоя
	srv := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Запускаем сервер в отдельной горутине
	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Ожидаем сигнал завершения работы (SIGINT, SIGTERM)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Printf("Shutdown signal received")

	// Выполняем graceful shutdown с таймаутом 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return
	}
	log.Printf("Server exited cleanly")
}
