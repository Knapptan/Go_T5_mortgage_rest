// Тесты middleware и routes
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

type MockHandler struct {
	mock.Mock
}

func (m *MockHandler) Execute(c *gin.Context) {
	m.Called(c)
}

func (m *MockHandler) GetCache(c *gin.Context) {
	m.Called(c)
}

func TestLoggingMiddleware(t *testing.T) {
	// Перехватываем вывод лога
	var logOutput bytes.Buffer
	originalOutput := log.Writer()
	log.SetOutput(&logOutput)
	defer log.SetOutput(originalOutput)

	// Создаем тестовый роутер с middleware
	router := gin.New()
	router.Use(LoggingMiddleware())
	router.GET("/test", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	// Выполняем запрос
	w := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/test", nil)
	router.ServeHTTP(w, req)

	// Проверяем вывод
	assert.Contains(t, logOutput.String(), "status_code: 200")
	assert.Contains(t, logOutput.String(), "duration: ")
}

func TestSetupRoutes(t *testing.T) {
	// Создаем мок
	mockHandler := new(MockHandler)

	// Настраиваем роутер
	router := gin.Default()
	SetupRoutes(router, mockHandler)

	// Проверяем зарегистрированные роуты
	routes := router.Routes()
	assert.Len(t, routes, 2)

	// Проверяем POST /execute
	assert.Equal(t, "POST", routes[0].Method)
	assert.Equal(t, "/execute", routes[0].Path)

	// Проверяем GET /cache
	assert.Equal(t, "GET", routes[1].Method)
	assert.Equal(t, "/cache", routes[1].Path)
}
