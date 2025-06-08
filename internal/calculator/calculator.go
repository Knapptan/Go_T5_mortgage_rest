// Package calculator содержит бизнес логику крелитного калькулятора
package calculator

import (
	"errors"
	"math"
	"time"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// Ошибки для возврата соответсвующего тела ответа
var (
	ErrNoProgramSelected   = errors.New("choose program")
	ErrMultiplePrograms    = errors.New("choose only 1 program")
	ErrInsufficientPayment = errors.New("the initial payment should be more")
)

// TODO попробовать сделать на децимал

// Рассчет запрашиваемых параметров кредита и запись в MortgageResponse
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
	annuity := calculateAnnuity(loanSum, rate, req.Months)
	totalPayment := annuity * float64(req.Months)
	overpayment := totalPayment - loanSum
	lasDate := time.Now().AddDate(0, req.Months, 0).Format(models.DateFormat)

	// Формирование ответа
	resp := models.MortgageResponse{
		Params: models.MortgageParams{
			ObjectCost:     req.ObjectCost,
			InitialPayment: req.InitialPayment,
			Months:         req.Months,
		},
		Program: req.Program,
		Aggregates: models.MortgageAggregates{
			Rate:            rate,
			LoanSum:         math.Round(loanSum*100) / 100,
			MonthlyPayment:  math.Round(annuity*100) / 100,
			Overpayment:     math.Round(overpayment*100) / 100,
			LastPaymentDate: lasDate,
		},
	}

	return resp, nil
}

// Рассчет размера ежемесячного аннуитетного платежа
func calculateAnnuity(loanSum, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100
	if monthlyRate == 0 { // Защита от деления на ноль
		return loanSum / float64(months)
	}
	discountFactor := math.Pow(1+monthlyRate, float64(months))
	annuity := loanSum * monthlyRate * discountFactor / (discountFactor - 1)
	return math.Round(annuity*100) / 100 // Округление до копеек
}
