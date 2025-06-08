// Package http содержит HTTP обработчики
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// Интерфейсы для зависимостей
type MortgageCalculator interface {
	Calculate(req models.MortgageRequest) (models.MortgageResponse, error)
}

type MortgageCache interface {
	Add(response models.MortgageResponse)
	GetAll() []models.MortgageInfoResponse
	IsEmpty() bool
}

// Структура с инъекцией зависимости
type Handler struct {
	calculator MortgageCalculator
	cache      MortgageCache
}

// Конструктр - принимает интерфейсы калькулятора и кэша
func NewHandler(calc MortgageCalculator, cache MortgageCache) *Handler {
	return &Handler{
		calculator: calc,
		cache:      cache,
	}
}

// Ручка обработки /execute, метод по указателю для изменнеия кэша, парсит тело запроса в MortgageRequest, вызывает логику и возвращает 200 с результатом в теле или 400 при ошибке c ответом в теле
func (h *Handler) Execute(c *gin.Context) {
	var req models.MortgageRequest
	// Парсинг Json из тела запроса
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	// Вызов логики, обработка логических ошибок
	resp, err := h.calculator.Calculate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Кэширование
	h.cache.Add(resp)

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// Ручка обработкии /cache возвращает все закешированые значения
func (h *Handler) GetCache(c *gin.Context) {
	if h.cache.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty cache"})
		return
	}

	c.JSON(http.StatusOK, h.cache.GetAll())
}
