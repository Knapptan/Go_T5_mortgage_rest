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

// NewCache - конструктор
func NewCache() *Cache {
	return &Cache{
		items:  make([]models.MortgageInfoResponse, 0),
		nextID: 0,
	}
}

// Add добавляет в кэш стурктуры MortgageResponse (с добавлением ID)
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

// GetAll возвращает копию среза items, чтобы внешние мутации не влияли на внутренний кеш.
func (c *Cache) GetAll() []models.MortgageInfoResponse {
	c.mu.RLock()
	defer c.mu.RUnlock()

	// Создаём новый срез нужной длины
	result := make([]models.MortgageInfoResponse, len(c.items))
	// Копируем все элементы
	copy(result, c.items)
	return result
}

// Get для получения записи по ID
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

// IsEmpty проверяет на пустоту
func (c *Cache) IsEmpty() bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return len(c.items) == 0
}
