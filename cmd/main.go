// Package main запускает HTTP-сервер Mortgage Calculator
package main

import (
	"log"
	"strconv"

	"github.com/Knapptan/Go_T5_mortgage_rest/config"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/cache"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/http"
	"github.com/gin-gonic/gin"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Инициализация компонентов
	mortgageCache := cache.New()
	handler := http.NewHandler(mortgageCache)

	// Настройка роутера
	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	http.SetupRoutes(router, handler)

	// Запуск сервера
	log.Printf("Starting server on port %d", cfg.Port)
	if err := router.Run(":" + strconv.Itoa(cfg.Port)); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
