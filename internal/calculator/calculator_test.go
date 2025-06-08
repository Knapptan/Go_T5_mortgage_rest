// Package calculator_test содержит модульные тесты основной логики пакета calculator.
package calculator_test

import (
	"testing"
	"time"

	. "github.com/Knapptan/Go_T5_mortgage_rest/internal/calculator"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestCalculate проверяет корректность расчётов для валидных входных данных с разными программами ипотеки.
func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		request  models.MortgageRequest
		expected models.MortgageResponse
		err      error
	}{
		{
			name: "valid salary program 0",
			request: models.MortgageRequest{
				ObjectCost:     5000000,
				InitialPayment: 1000000,
				Months:         240,
				Program: struct {
					Salary   bool `json:"salary"`
					Military bool `json:"military"`
					Base     bool `json:"base"`
				}{Salary: true},
			},
			expected: models.MortgageResponse{
				Aggregates: struct {
					Rate            float64 `json:"rate"`
					LoanSum         float64 `json:"loan_sum"`
					MonthlyPayment  float64 `json:"monthly_payment"`
					Overpayment     float64 `json:"overpayment"`
					LastPaymentDate string  `json:"last_payment_date"`
				}{
					Rate:            8,
					LoanSum:         4000000,
					MonthlyPayment:  33458,
					Overpayment:     4029920,
					LastPaymentDate: time.Now().AddDate(0, 240, 0).Format(models.DateFormat),
				},
			},
		},
		{
			name: "valid salary program 1",
			request: models.MortgageRequest{
				ObjectCost:     8000000,
				InitialPayment: 2000000,
				Months:         200,
				Program: struct {
					Salary   bool `json:"salary"`
					Military bool `json:"military"`
					Base     bool `json:"base"`
				}{Military: true},
			},
			expected: models.MortgageResponse{
				Aggregates: struct {
					Rate            float64 `json:"rate"`
					LoanSum         float64 `json:"loan_sum"`
					MonthlyPayment  float64 `json:"monthly_payment"`
					Overpayment     float64 `json:"overpayment"`
					LastPaymentDate string  `json:"last_payment_date"`
				}{
					Rate:            9,
					LoanSum:         6000000,
					MonthlyPayment:  58019,
					Overpayment:     5603800,
					LastPaymentDate: time.Now().AddDate(0, 200, 0).Format(models.DateFormat),
				},
			},
		},
		{
			name: "valid salary program 2",
			request: models.MortgageRequest{
				ObjectCost:     12000000,
				InitialPayment: 3000000,
				Months:         120,
				Program: struct {
					Salary   bool `json:"salary"`
					Military bool `json:"military"`
					Base     bool `json:"base"`
				}{Base: true},
			},
			expected: models.MortgageResponse{
				Aggregates: struct {
					Rate            float64 `json:"rate"`
					LoanSum         float64 `json:"loan_sum"`
					MonthlyPayment  float64 `json:"monthly_payment"`
					Overpayment     float64 `json:"overpayment"`
					LastPaymentDate string  `json:"last_payment_date"`
				}{
					Rate:            10,
					LoanSum:         9000000,
					MonthlyPayment:  118936,
					Overpayment:     5272320,
					LastPaymentDate: time.Now().AddDate(0, 120, 0).Format(models.DateFormat),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := Calculate(tt.request)

			if tt.err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.err, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.expected.Aggregates.Rate, resp.Aggregates.Rate)
			assert.InDelta(t, tt.expected.Aggregates.MonthlyPayment, resp.Aggregates.MonthlyPayment, 1.0)
			assert.Equal(t, tt.expected.Aggregates.LastPaymentDate, resp.Aggregates.LastPaymentDate)
		})
	}
}

// TestCalculate_InvalidParameters проверяет ошибки валидации при некорректных входных данных.
func TestCalculate_InvalidParameters(t *testing.T) {
	tests := []struct {
		name    string
		req     models.MortgageRequest
		wantErr error
	}{
		{
			name: "zero program",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 500000,
				Months:         120,
				Program:        models.MortgageProgram{},
			},
			wantErr: ErrNoProgramSelected,
		},
		{
			name: "more then one program",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 500000,
				Months:         120,
				Program:        models.MortgageProgram{Base: true, Military: true},
			},
			wantErr: ErrMultiplePrograms,
		},
		{
			name: "zero initial payment",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 0,
				Months:         120,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrInsufficientPayment,
		},
		{
			name: "negative object cost",
			req: models.MortgageRequest{
				ObjectCost:     -1000000,
				InitialPayment: 200000,
				Months:         120,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrInvalidParameters,
		},
		{
			name: "negative initial payment",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: -1,
				Months:         120,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrInvalidParameters,
		},
		{
			name: "initial payment exceed object cost",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 1000001,
				Months:         120,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrExcessivePayment,
		},
		{
			name: "invalid duration",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 200000,
				Months:         0,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrInvalidDuration,
		},
		{
			name: "initial payment exceeds cost",
			req: models.MortgageRequest{
				ObjectCost:     1000000,
				InitialPayment: 2000000,
				Months:         120,
				Program:        models.MortgageProgram{Base: true},
			},
			wantErr: ErrExcessivePayment,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := Calculate(tt.req)
			if tt.wantErr == nil {
				require.NoError(t, err)
			} else {
				require.ErrorContains(t, err, tt.wantErr.Error())
			}
		})
	}
}

// TestNewService проверяет, что NewService возвращает непустой экземпляр Service.
func TestNewService(t *testing.T) {
	s := NewService()
	assert.NotNil(t, s)
}

// TestNewService проверяет, что NewService возвращает непустой экземпляр Service.
func TestService_Calculate_WrapsCalculate(t *testing.T) {
	svc := NewService()

	// Подготовим "валидный" запрос, который Calculate успешно обработает
	req := models.MortgageRequest{
		ObjectCost:     1000,
		InitialPayment: 200,
		Months:         12,
		Program:        models.MortgageProgram{Salary: true},
	}

	// Получаем ожидаемый результат, вызывая напрямую Calculate
	wantResp, wantErr := Calculate(req)

	// Вызываем метод обёртки
	gotResp, gotErr := svc.Calculate(req)

	// Сравниваем оба результата
	assert.Equal(t, wantErr, gotErr, "Service.Calculate should return same error as Calculate")
	assert.Equal(t, wantResp, gotResp, "Service.Calculate should return same response as Calculate")
}
