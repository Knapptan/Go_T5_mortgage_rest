// Тесты пакета http
package http_test

import (
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/mock"
)

// MockCache для тестирования
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

// MockCalculator для тестирования
type MockCalculator struct {
	mock.Mock
}

func (m *MockCalculator) Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	args := m.Called(req)
	return args.Get(0).(models.MortgageResponse), args.Error(1)
}
