// Тесты пакета cache
package cache_test

import (
	"testing"

	. "github.com/Knapptan/Go_T5_mortgage_rest/internal/cache"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/assert"
)

// Тест конструктора и начального состояния
func TestNewCache_InitialState(t *testing.T) {
	c := NewCache()
	assert.NotNil(t, c)
	assert.True(t, c.IsEmpty())
	assert.Equal(t, 0, len(c.GetAll()))
}

// Тест Add: проставляет ID начиная с 0, увеличивает nextID и сохраняет данные
func TestCache_AddAssignsIncrementalIDs(t *testing.T) {
	c := NewCache()

	// Два ответа с пустыми полями
	resp1 := models.MortgageResponse{}
	resp2 := models.MortgageResponse{}

	c.Add(resp1)
	c.Add(resp2)

	all := c.GetAll()
	assert.Len(t, all, 2)

	// Первому должна быть присвоена ID=0, второму ID=1
	assert.Equal(t, 0, all[0].ID)
	assert.Equal(t, 1, all[1].ID)
}

// Тест GetAll: возвращает копию среза, не ссылку на внутренний
func TestCache_GetAll_ReturnsCopy(t *testing.T) {
	c := NewCache()
	resp := models.MortgageResponse{}
	c.Add(resp)

	all1 := c.GetAll()
	all1[0].ID = 999 // мутируем возвращённый срез

	all2 := c.GetAll()
	// Внутренний срез не изменился
	assert.Equal(t, 0, all2[0].ID)
}

// Тест Get: существующий ID
func TestCache_GetExisting(t *testing.T) {
	c := NewCache()
	resp := models.MortgageResponse{}
	c.Add(resp) // ID=0

	item, found := c.Get(0)
	assert.True(t, found)
	assert.Equal(t, 0, item.ID)
}

// Тест Get: несуществующий ID
func TestCache_GetNotFound(t *testing.T) {
	c := NewCache()
	// Ни одного Add — ID=0 не существует
	_, found := c.Get(0)
	assert.False(t, found)

	// Добавим один, с ID=0
	c.Add(models.MortgageResponse{})
	// Запросим ID=1 — тоже не должно найти
	_, found = c.Get(1)
	assert.False(t, found)
}

// Тест IsEmpty после операций
func TestCache_IsEmptyBehavior(t *testing.T) {
	c := NewCache()
	assert.True(t, c.IsEmpty())

	c.Add(models.MortgageResponse{})
	assert.False(t, c.IsEmpty())

	// После добавления одного элемента GetAll непустой
	assert.NotEmpty(t, c.GetAll())
}
