// Package main запускает HTTP-сервер Mortgage Calculator.
package main

import (
	"log"

	"github.com/Knapptan/Go_T5_mortgage_rest/config"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	log.Println(cfg)
}
