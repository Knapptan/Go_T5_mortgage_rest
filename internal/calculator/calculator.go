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
	ErrNoProgramSelected   = errors.New("choose program")                                   // ни одна ипотечная программа не выбрана
	ErrMultiplePrograms    = errors.New("choose only 1 program")                            // выбрано более одной программы
	ErrInsufficientPayment = errors.New("the initial payment should be more")               // взнос меньше 20% стоимости
	ErrInvalidParameters   = errors.New("object cost and initial payment must be positive") // некорректные параметры стоимости или взноса
	ErrInvalidDuration     = errors.New("loan duration must be at least 1 month")           // срок кредита меньше 1 месяца
	ErrExcessivePayment    = errors.New("initial payment cannot exceed object cost")        // взнос превышает стоимость
)

// Calculate делегирует вызов функции Calculate на уровне пакета и реализует метод интерфейса.
func (s *Service) Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	return Calculate(req)
}

// Calculate выполняет полный цикл: валидацию параметров и расчёт ипотечных агрегатов.
func Calculate(req models.MortgageRequest) (models.MortgageResponse, error) {
	// 1. Валидация базовых параметров
	if err := validateRequest(req); err != nil {
		return models.MortgageResponse{}, err
	}
	// 2. Валидация выбора программы
	if err := validateProgram(req.Program); err != nil {
		return models.MortgageResponse{}, err
	}
	// 3. Валидация первоначального взноса
	if err := validatePayment(req); err != nil {
		return models.MortgageResponse{}, err
	}

	// 4. Определяем годовую ставку в зависимости от программы
	rate := selectRate(req.Program)

	// 5. Вычисляем сумму кредита после взноса
	loanSum := req.ObjectCost - req.InitialPayment
	// 6. Рассчитываем аннуитетный платёж
	annuity := calculateAnnuity(loanSum, rate, req.Months)
	// 7. Формируем дату последнего платежа (год и месяц от now + Months)
	lastDate := time.Now().AddDate(0, req.Months, 0).Format(models.DateFormat)

	// 8. Собираем окончательный ответ
	resp := models.MortgageResponse{
		Params: models.MortgageParams{
			ObjectCost:     req.ObjectCost,     // исходная стоимость
			InitialPayment: req.InitialPayment, // первоначальный взнос
			Months:         req.Months,         // срок кредита
		},
		Program: req.Program, // выбранная ипотечная программа
		Aggregates: models.MortgageAggregates{
			Rate:            rate,                                                        // годовая ставка
			LoanSum:         math.Round(loanSum*100) / 100,                               // сумма кредита с округлением
			MonthlyPayment:  math.Round(annuity*100) / 100,                               // ежемесячный платёж
			Overpayment:     math.Round((annuity*float64(req.Months)-loanSum)*100) / 100, // переплата
			LastPaymentDate: lastDate,                                                    // дата последнего платежа
		},
	}
	return resp, nil // возвращаем заполненный ответ без ошибки
}

// validateRequest проверяет корректность основных числовых параметров.
func validateRequest(req models.MortgageRequest) error {
	if req.ObjectCost <= 0 || req.InitialPayment < 0 {
		return ErrInvalidParameters // неверные параметры
	}
	if req.Months < 1 {
		return ErrInvalidDuration // срок меньше 1 месяца
	}
	return nil
}

// validateProgram проверяет, что выбрана ровно одна ипотечная программа.
func validateProgram(p models.MortgageProgram) error {
	count := 0
	if p.Salary {
		count++
	}
	if p.Military {
		count++
	}
	if p.Base {
		count++
	}
	switch {
	case count == 0:
		return ErrNoProgramSelected
	case count > 1:
		return ErrMultiplePrograms
	}
	return nil
}

// validatePayment проверяет, что первоначальный взнос в допустимых пределах.
func validatePayment(req models.MortgageRequest) error {
	if req.InitialPayment > req.ObjectCost {
		return ErrExcessivePayment
	}
	if req.InitialPayment < req.ObjectCost*0.2 {
		return ErrInsufficientPayment
	}
	return nil
}

// selectRate возвращает годовую ставку в зависимости от выбранной программы.
func selectRate(p models.MortgageProgram) float64 {
	switch {
	case p.Salary:
		return 8.0 // корпоративная
	case p.Military:
		return 9.0 // военная
	default:
		return 10.0 // базовая
	}
}

// calculateAnnuity рассчитывает ежемесячный аннуитетный платёж и округляет до копеек.
func calculateAnnuity(loanSum, annualRate float64, months int) float64 {
	monthlyRate := annualRate / 12 / 100 // переводим годовую ставку в месячную (доля от 1)
	if monthlyRate == 0 {
		return loanSum / float64(months) // защита от деления на ноль
	}
	discountFactor := math.Pow(1+monthlyRate, float64(months))
	annuity := loanSum * monthlyRate * discountFactor / (discountFactor - 1)
	return math.Round(annuity*100) / 100 // округляем до двух знаков
}
