// Package http предоставляет HTTP-обработчики для работы с ипотечными расчётами и кэшем.
package http

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// MortgageCalculator определяет интерфейс для расчёта ипотеки.
type MortgageCalculator interface {
	Calculate(req models.MortgageRequest) (models.MortgageResponse, error)
}

// MortgageCache определяет интерфейс для работы с кэшем ипотечных расчётов.
type MortgageCache interface {
	Add(response models.MortgageResponse)
	GetAll() []models.MortgageInfoResponse
	IsEmpty() bool
}

// Handler обрабатывает HTTP-запросы и использует зависимости калькулятора и кэша.
type Handler struct {
	calculator MortgageCalculator
	cache      MortgageCache
}

// NewHandler создаёт новый Handler с переданными зависимостями калькулятора и кэша.
func NewHandler(calc MortgageCalculator, cache MortgageCache) *Handler {
	if calc == nil {
		panic("MortgageCalculator dependency is nil")
	}
	if cache == nil {
		panic("MortgageCache dependency is nil")
	}
	return &Handler{
		calculator: calc,
		cache:      cache,
	}
}

// Execute обрабатывает POST-запрос на /execute, выполняет расчёт ипотеки и возвращает результат.
func (h *Handler) Execute(c *gin.Context) {
	if c.Request.Method != http.MethodPost {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
		return
	}

	var req models.MortgageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
		return
	}

	resp, err := h.calculator.Calculate(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	h.cache.Add(resp)

	c.JSON(http.StatusOK, gin.H{"result": resp})
}

// GetCache обрабатывает GET-запрос на /cache и возвращает все закешированные значения.
func (h *Handler) GetCache(c *gin.Context) {
	if c.Request.Method != http.MethodGet {
		c.JSON(http.StatusMethodNotAllowed, gin.H{"error": "method not allowed"})
		return
	}

	if h.cache.IsEmpty() {
		c.JSON(http.StatusBadRequest, gin.H{"error": "empty cache"})
		return
	}

	c.JSON(http.StatusOK, h.cache.GetAll())
}
