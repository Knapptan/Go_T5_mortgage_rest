// Тесты handler-ов
package http_test

import (
	"github.com/gin-gonic/gin"
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
