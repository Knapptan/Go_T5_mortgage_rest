// Package cache_test содержит модульные тесты для пакета cache.
package cache_test

import (
	"testing"

	. "github.com/Knapptan/Go_T5_mortgage_rest/internal/cache"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/assert"
)

// TestNewCache_InitialState проверяет корректность инициализации нового кэша.
func TestNewCache_InitialState(t *testing.T) {
	c := NewCache()
	assert.NotNil(t, c)
	assert.True(t, c.IsEmpty())
	assert.Equal(t, 0, len(c.GetAll()))
}

// TestCache_AddAssignsIncrementalIDs проверяет, что Add присваивает элементы с автоинкрементным ID.
func TestCache_AddAssignsIncrementalIDs(t *testing.T) {
	c := NewCache()

	resp1 := models.MortgageResponse{}
	resp2 := models.MortgageResponse{}

	c.Add(resp1)
	c.Add(resp2)

	all := c.GetAll()
	assert.Len(t, all, 2)
	assert.Equal(t, 0, all[0].ID)
	assert.Equal(t, 1, all[1].ID)
}

// TestCache_GetAll_ReturnsCopy проверяет, что GetAll возвращает копию среза, а не ссылку на внутренние данные.
func TestCache_GetAll_ReturnsCopy(t *testing.T) {
	c := NewCache()
	resp := models.MortgageResponse{}
	c.Add(resp)

	all1 := c.GetAll()
	all1[0].ID = 999

	all2 := c.GetAll()
	assert.Equal(t, 0, all2[0].ID)
}

// TestCache_GetExisting проверяет получение существующего элемента по ID.
func TestCache_GetExisting(t *testing.T) {
	c := NewCache()
	resp := models.MortgageResponse{}
	c.Add(resp)

	item, found := c.Get(0)
	assert.True(t, found)
	assert.Equal(t, 0, item.ID)
}

// TestCache_GetNotFound проверяет поведение Get при запросе несуществующего ID.
func TestCache_GetNotFound(t *testing.T) {
	c := NewCache()

	_, found := c.Get(0)
	assert.False(t, found)

	c.Add(models.MortgageResponse{})
	_, found = c.Get(1)
	assert.False(t, found)
}

// TestCache_IsEmptyBehavior проверяет корректность работы метода IsEmpty.
func TestCache_IsEmptyBehavior(t *testing.T) {
	c := NewCache()
	assert.True(t, c.IsEmpty())

	c.Add(models.MortgageResponse{})
	assert.False(t, c.IsEmpty())
	assert.NotEmpty(t, c.GetAll())
}
