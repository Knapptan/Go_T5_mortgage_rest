package http_test

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Knapptan/Go_T5_mortgage_rest/internal/http"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockHandler реализует интерфейс MortgageCache для целей тестирования.
type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Execute(c *gin.Context) {
	m.Called(c)
}

func (m *MockHandler) GetCache(c *gin.Context) {
	m.Called(c)
}

// TestLoggingMiddleware проверяет, что LoggingMiddleware корректно логирует успешные HTTP-запросы.
func TestLoggingMiddleware(t *testing.T) {
	// Перенаправляем вывод логов в буфер для проверки
	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logOutput)
	defer log.SetOutput(originalOutput)

	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	assert.Contains(t, logOutput.String(), "status_code: 200")
	assert.Contains(t, logOutput.String(), "duration: ")
}

// TestLoggingMiddleware_ErrorStatus проверяет логирование при ответе с ошибочным статусом.
func TestLoggingMiddleware_ErrorStatus(t *testing.T) {
	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logOutput)
	defer log.SetOutput(originalOutput)

	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/test-error", func(c *gin.Context) {
		c.Status(http.StatusInternalServerError)
	})

	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test-error", nil)
	router.ServeHTTP(w, req)

	assert.Contains(t, logOutput.String(), "status_code: 500")
	assert.Contains(t, logOutput.String(), "duration: ")
}

// TestSetupRoutes проверяет, что SetupRoutes регистрирует корректные маршруты.
func TestSetupRoutes(t *testing.T) {
	mockHandler := new(MockHandler)

	mockHandler.On("Execute", mock.Anything).Return()
	mockHandler.On("GetCache", mock.Anything).Return()

	router := gin.Default()

	SetupRoutes(router, mockHandler)

	routes := router.Routes()
	assert.Len(t, routes, 2)

	routeMap := make(map[string]string)
	for _, route := range routes {
		routeMap[route.Path+"::"+route.Method] = route.Path
	}

	assert.Contains(t, routeMap, "/execute::POST")
	assert.Contains(t, routeMap, "/cache::GET")
}
