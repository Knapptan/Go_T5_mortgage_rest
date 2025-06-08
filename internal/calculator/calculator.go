// Package calculator содержит бизнес-логику ипотечного калькулятора.
package calculator

import (
	"errors"
	"math"
	"time"

	"github.com/Knapptan/Go_T5_mortgage_rest/internal/models"
)

// Service реализует интерфейс ипотечного калькулятора.
type Service struct{}

// NewService возвращает новый экземпляр Service.
func NewService() *Service {
	return &Service{}
}

// Предопределённые ошибки, возвращаемые при валидации входных параметров.
var (
	ErrNoProgramSelected   = errors.New("choose program")                                   // не выбрана ни одна программа
	ErrMultiplePrograms    = errors.New("choose only 1 program")                            // выбрано несколько программ
	ErrInsufficientPayment = errors.New("the initial payment should be more")               // первоначальный взнос меньше требуемого минимума
	ErrInvalidParameters   = errors.New("object cost and initial payment must be positive") // некорректные значения стоимости или взноса
	ErrInvalidDuration     = errors.New("loan duration must be at least 1 month")           // срок кредита меньше одного месяца
	ErrExcessivePayment    = errors.New("initial payment cannot exceed object cost")        // взнос превышает стоимость объекта
)

// Calculate делегирует вычисление функции Calculate и реализует интерфейс MortgageCalculator.
func (s *Service) Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	return Calculate(req)
}

// Calculate выполняет валидацию и расчёт параметров кредита на основе входного запроса.
func Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	// Проверка входных параметров
	if req.ObjectCost <= 0 || req.InitialPayment < 0 {
		return models.MortgageResponse{}, ErrInvalidParameters
	}
	if req.Months < 1 {
		return models.MortgageResponse{}, ErrInvalidDuration
	}

	// Проверка выбора программы
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

	// Проверка валидности взноса
	if req.InitialPayment > req.ObjectCost {
		return models.MortgageResponse{}, ErrExcessivePayment
	}
	minPayment := req.ObjectCost * 0.2
	if req.InitialPayment < minPayment {
		return models.MortgageResponse{}, ErrInsufficientPayment
	}

	// Выбор процентной ставки в зависимости от программы
	rate := 0.0
	switch {
	case req.Program.Salary:
		rate = 8.0
	case req.Program.Military:
		rate = 9.0
	case req.Program.Base:
		rate = 10.0
	}

	// Вычисление итоговых значений
	loanSum := req.ObjectCost - req.InitialPayment
	annuity := calculateAnnuity(loanSum, rate, req.Months)
	totalPayment := annuity * float64(req.Months)
	overpayment := totalPayment - loanSum
	lastDate := time.Now().AddDate(0, req.Months, 0).Format(models.DateFormat)

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
			LastPaymentDate: lastDate,
		},
	}

	return resp, nil
}

// calculateAnnuity рассчитывает размер ежемесячного аннуитетного платежа.
func calculateAnnuity(loanSum, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100
	if monthlyRate == 0 {
		return loanSum / float64(months)
	}
	discountFactor := math.Pow(1+monthlyRate, float64(months))
	annuity := loanSum * monthlyRate * discountFactor / (discountFactor - 1)
	return math.Round(annuity*100) / 100
}
