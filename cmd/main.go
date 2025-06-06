// Package main запускает HTTP-сервер Mortgage Calculator
package main

import (
	"log"

	"github.com/Knapptan/Go_T5_mortgage_rest/config"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/cache"
)

func main() {
	// Загрузка конфигурации
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg) // TODO убрать 

	mortgageCache:= cache.New() 

	log.Println(mortgageCache) // TODO убрать 
}
