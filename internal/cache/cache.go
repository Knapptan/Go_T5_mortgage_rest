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
func NewCache() *Cache {
	return &Cache{
		items:  make([]models.MortgageInfoResponse, 0),
		nextID: 0,
	}
}

// Метод добавления в кэш стурктуры MortgageResponse (с добавлением ID)
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
	c.set(item)
}

// Метод добавления в массив кэша (приватный, вынесен для тестов)
func (c *Cache) set(item models.MortgageInfoResponse) {
	c.items = append(c.items, item)
}

// Метод получения всех записей возвращает items
func (c *Cache) GetAll() []models.MortgageInfoResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.items
}

// Метод получения записи по ID
// Если элемент найден, возвращает (item, true).
// Если нет — (zero, false).
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

// Метод проверки на пустоту
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}
