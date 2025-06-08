// Тесты пакета http
package http_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	. "github.com/Knapptan/Go_T5_mortgage_rest/internal/http"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockCache реализует MortgageCache
type MockCache struct {
	mock.Mock
}

func (m *MockCache) Add(response models.MortgageResponse) {
	m.Called(response)
}

func (m *MockCache) GetAll() []models.MortgageInfoResponse {
	args := m.Called()
	return args.Get(0).([]models.MortgageInfoResponse)
}

func (m *MockCache) IsEmpty() bool {
	args := m.Called()
	return args.Bool(0)
}

// MockCalculator реализует MortgageCalculator
type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.MortgageResponse), args.Error(1)
}

// Вспомогательная функция для преобразования в JSON
func mustJSON(v interface{}) []byte {
	data, _ := json.Marshal(v)
	return data
}

func TestExecuteHandler_Success(t *testing.T) {
	// Инициализация моков
	mockCalc := new(MockCalculator)
	mockCache := new(MockCache)
	handler := NewHandler(mockCalc, mockCache)

	// Тестовые данные
	req := models.MortgageRequest{
		ObjectCost:     5000000,
		InitialPayment: 1000000,
		Months:         240,
		Program:        models.MortgageProgram{Salary: true},
	}

	resp := models.MortgageResponse{
		Aggregates: models.MortgageAggregates{
			MonthlyPayment: 33458.33,
		},
	}

	// Настройка ожиданий
	mockCalc.On("Calculate", req).Return(resp, nil)
	mockCache.On("Add", resp).Return()

	// Создание тестового контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("POST", "/execute", bytes.NewReader(mustJSON(req)))
	c.Request.Header.Set("Content-Type", "application/json")

	handler.Execute(c)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)

	var response struct {
		Result models.MortgageResponse `json:"result"`
	}
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Equal(t, 33458.33, response.Result.Aggregates.MonthlyPayment)
	mockCalc.AssertExpectations(t)
	mockCache.AssertExpectations(t)
}

func TestGetCacheHandler_Success(t *testing.T) {
	mockCalc := new(MockCalculator)
	mockCache := new(MockCache)
	handler := NewHandler(mockCalc, mockCache)

	// Ожидаемые данные кэша
	expectedItems := []models.MortgageInfoResponse{
		{
			ID: 1,
			Aggregates: models.MortgageAggregates{
				MonthlyPayment: 33458.33,
			},
		},
	}

	// Настройка ожиданий
	mockCache.On("IsEmpty").Return(false)
	mockCache.On("GetAll").Return(expectedItems)

	// Создание тестового контекста
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest("GET", "/cache", nil)

	// Вызов обработчика
	handler.GetCache(c)

	// Проверки
	assert.Equal(t, http.StatusOK, w.Code)

	var response []models.MortgageInfoResponse
	json.Unmarshal(w.Body.Bytes(), &response)

	assert.Len(t, response, 1)
	assert.Equal(t, 33458.33, response[0].Aggregates.MonthlyPayment)
	mockCache.AssertExpectations(t)
}
