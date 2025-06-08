// Package calculator_test содержит тесты осноной логики сервиса
package calculator_test

import (
	"testing"
	"time"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/calculator"
	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
	"github.com/stretchr/testify/assert"
)

func TestCalculate(t *testing.T) {
	tests := []struct {
		name     string
		request  models.MortgageRequest
		expected models.MortgageResponse
		err      error
	}{
		{
			name: "Valid salary program 0",
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
			name: "Valid salary program 1",
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
			name: "Valid salary program 2",
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
			resp, err := calculator.Calculate(tt.request)

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
