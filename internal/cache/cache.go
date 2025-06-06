// Package cache это реализация кэша
package cache

import (
	"sync"

	"github.com/Knapptan/Go_T5_mortgage_rest/pkg/models"
)

// Структура кэша с рв-мютексом во избежание data races
type Cache struct {
	mu     sync.RWMutex
	items  []models.CacheItem
	nextID int
}

// Конструктор
func New() *Cache {
	return &Cache{
		items:  make([]models.CacheItem, 0),
		nextID: 1,
	}
}

// Метод добавления
func (c *Cache) Add(item models.CacheItem) {
	c.mu.Lock()
	defer c.mu.Unlock()

	item.ID = c.nextID
	c.nextID++
	c.items = append(c.items, item)
}

// Метод получения всех записей
func (c *Cache) GetAll() []models.CacheItem {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make([]models.CacheItem, len(c.items))
	copy(result, c.items)
	return result
}

// Метод проверки на пустоту
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}
