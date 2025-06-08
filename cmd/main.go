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
	myhttp "github.com/Knapptan/Go_T5_mortgage_rest/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация компонентов
	mortgageCache := cache.NewCache()
	handler := myhttp.NewHandler(mortgageCache)

	// Настройка роутера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	myhttp.SetupRoutes(router, handler)

	// Создаем HTTP-сервер с таймаутами
	srv := &http.Server{
		Addr:    ":" + strconv.Itoa(cfg.Port),
		Handler: router,
		// Рекомендуемые таймауты для production
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  15 * time.Second,
	}

	// Запуск сервера в горутине
	go func() {
		log.Printf("Starting server on port %d", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("server error: %v", err)
		}
	}()

	// Канал для сигналов ОС
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	// Блокируем, пока не получим сигнал
	<-quit
	log.Printf("Shutdown signal received")

	// Создаем контекст с таймаутом для graceful shutdown 5 секунд
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Останавливаем сервер
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown: %v", err)
		return
	}
	log.Printf("Server exited cleanly")
}
