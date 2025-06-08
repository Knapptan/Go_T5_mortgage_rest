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
	// Создаем мок обработчика
	mockHandler := new(MockHandler)

	// Настраиваем ожидания вызовов
	mockHandler.On("Execute", mock.Anything).Return()
	mockHandler.On("GetCache", mock.Anything).Return()

	// Создаем роутер
	router := gin.Default()

	// Настраиваем маршруты
	SetupRoutes(router, mockHandler)

	// Проверяем зарегистрированные роуты
	routes := router.Routes()
	assert.Len(t, routes, 2)

	// Собираем маршруты в map для удобной проверки
	routeMap := make(map[string]string)
	for _, route := range routes {
		routeMap[route.Path+"::"+route.Method] = route.Path
	}

	// Проверяем наличие ожидаемых маршрутов
	assert.Contains(t, routeMap, "/execute::POST")
	assert.Contains(t, routeMap, "/cache::GET")
}
