// Package cache предоставляет потокобезопасную реализацию кэша для ипотечных расчётов.
package cache

import (
	"sync"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// Cache представляет потокобезопасный кэш с RWMutex для хранения ипотечных расчётов.
type Cache struct {
	items  []models.MortgageInfoResponse
	nextID int
	mu     sync.RWMutex
}

// NewCache создаёт и возвращает новый экземпляр Cache.
func NewCache() *Cache {
	return &Cache{
		items:  make([]models.MortgageInfoResponse, 0),
		nextID: 0,
	}
}

// Add добавляет новый элемент в кэш на основе MortgageResponse, присваивая уникальный ID.
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

// GetAll возвращает копию всех элементов кэша, чтобы избежать внешних изменений данных.
func (c *Cache) GetAll() []models.MortgageInfoResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]models.MortgageInfoResponse, len(c.items))
	copy(result, c.items)
	return result
}

// Get возвращает элемент по заданному ID.
// Если элемент найден, возвращает его и true; иначе — zero-значение и false.
func (c *Cache) Get(id int) (models.MortgageInfoResponse, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()

	for _, item := range c.items {
		if item.ID == id {
			return item, true
		}
	}

	return models.MortgageInfoResponse{}, false
}

// IsEmpty возвращает true, если кэш не содержит элементов.
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}
