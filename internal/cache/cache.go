// Package cache это реализация кэша
package cache

import (
	"sync"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// Структура кэша с рв-мютексом во избежание data races
type Cache struct {
	mu     sync.RWMutex
	items  []models.MortgageInfoResponse
	nextID int
}

// Конструктор
func New() *Cache {
	return &Cache{
		items:  make([]models.MortgageInfoResponse, 0),
		nextID: 1,
	}
}

// Метод добавления
func (c *Cache) Add(response models.MortgageResponse) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item := models.MortgageInfoResponse{
		ID:         c.nextID,
		Params:     response.Params,
		Program:    response.Program,
		Aggregates: response.Aggregates,
	}

	c.nextID++
	c.items = append(c.items, item)
}

// Метод получения всех записей
func (c *Cache) GetAll() []models.MortgageInfoResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

// Метод проверки на пустоту
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}
