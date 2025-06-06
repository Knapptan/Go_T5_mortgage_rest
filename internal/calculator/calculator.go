// Package calculator содержит бизнес логику крелитного калькулятора
package calculator

import (
	"errors"
	"math"
	"time"

	"github.com/Knapptan/Go_T5_mortgage_rest/pkg/models"
)

// Ошибки для возврата соответсвующего тела ответа
var (
	ErrNoProgramSelected   = errors.New("choose program")
	ErrMultiplePrograms    = errors.New("choose only 1 program")
	ErrInsufficientPayment = errors.New("the initial payment should be more")
)

func Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	// Валидация программы кредита
	programs := 0
	if req.Program.Salary {
		programs++
	}
	if req.Program.Military {
		programs++
	}
	if req.Program.Base {
		programs++
	}

	switch {
	case programs == 0:
		return models.MortgageResponse{}, ErrNoProgramSelected
	case programs > 1:
		return models.MortgageResponse{}, ErrMultiplePrograms
	}

	// Валидация взноса
	minPayment := req.ObjectCost * 0.2
	if req.InitialPayment < minPayment {
		return models.MortgageResponse{}, ErrInsufficientPayment
	}

	// Определение ставки
	rate := 0.0
	switch {
	case req.Program.Salary:
		rate = 8.0
	case req.Program.Military:
		rate = 9.0
	case req.Program.Base:
		rate = 10.0
	}

	// Расчет параметров
	loanSum := req.ObjectCost - req.InitialPayment
	monthlyRate := rate / 100 / 12
	annuity := calculateAnnuity(loanSum, monthlyRate, req.Months)
	totalPayment := annuity * float64(req.Months)
	overpayment := totalPayment - loanSum
	lasDate := time.Now().AddDate(0, req.Months, 0).Format("2006-01-02")

	// Формирование ответа
	resp := models.MortgageResponse{
		Params: struct {
			ObjectCost     float64 `json:"object_cost"`
			InitialPayment float64 `json:"initial_payment"`
			Months         int     `json:"months"`
		}{
			ObjectCost:     req.ObjectCost,
			InitialPayment: req.InitialPayment,
			Months:         req.Months,
		},
		Program: req.Program,
		Aggregates: struct {
			Rate            float64 `json:"rate"`
			LoanSum         float64 `json:"loan_sum"`
			MonthlyPayment  float64 `json:"monthly_payment"`
			Overpayment     float64 `json:"overpayment"`
			LastPaymentDate string  `json:"last_payment_date"`
		}{
			Rate:            rate,
			LoanSum:         math.Round(loanSum*100) / 100,
			MonthlyPayment:  math.Round(annuity*100) / 100,
			Overpayment:     math.Round(overpayment*100) / 100,
			LastPaymentDate: lasDate,
		},
	}
	return resp, nil
}

func calculateAnnuity(loanSum, monthlyRate float64, months int) float64 {
	return 0.0 // TODO logic
}
