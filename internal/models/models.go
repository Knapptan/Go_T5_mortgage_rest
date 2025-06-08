// Package models содержит структуры, описывающие входные и выходные данные ипотечного сервиса.
package models

// MortgageParams описывает основные параметры кредита:
// стоимость объекта, первоначальный взнос и срок в месяцах.
type MortgageParams struct {
	ObjectCost     float64 `json:"object_cost"`     // Стоимость объекта недвижимости.
	InitialPayment float64 `json:"initial_payment"` // Первоначальный взнос.
	Months         int     `json:"months"`          // Срок кредита в месяцах.
}

// MortgageProgram задаёт доступные ипотечные программы:
// корпоративную (Salary), военную (Military) и базовую (Base).
type MortgageProgram struct {
	Salary   bool `json:"salary"`   // Корпоративная ипотека.
	Military bool `json:"military"` // Военная ипотека.
	Base     bool `json:"base"`     // Базовая ипотечная программа.
}

// MortgageAggregates содержит рассчитанные значения по кредиту:
// процентную ставку, сумму кредита, размер ежемесячного платежа,
// общую переплату и дату последнего платежа.
type MortgageAggregates struct {
	LastPaymentDate string  `json:"last_payment_date"` // Дата последнего платежа в формате YYYY-MM-DD.
	Rate            float64 `json:"rate"`              // Годовая процентная ставка (в процентах).
	LoanSum         float64 `json:"loan_sum"`          // Общая сумма кредита после вычета первоначального взноса.
	MonthlyPayment  float64 `json:"monthly_payment"`   // Размер аннуитетного ежемесячного платежа.
	Overpayment     float64 `json:"overpayment"`       // Общая переплата за весь срок кредита.
}

// MortgageInfoResponse представляет одну запись в кэше:
// уникальный ID, параметры запроса и рассчитанные агрегаты.
type MortgageInfoResponse struct {
	Aggregates MortgageAggregates `json:"aggregates"` // Полученные агрегированные результаты.
	Params     MortgageParams     `json:"params"`     // Входные параметры кредита.
	Program    MortgageProgram    `json:"program"`    // Выбранная ипотечная программа.
	ID         int                // Уникальный идентификатор расчёта.
}

// MortgageRequest описывает JSON-запрос к сервису расчёта ипотечного кредита.
type MortgageRequest struct {
	ObjectCost     float64         `json:"object_cost"`     // Стоимость объекта недвижимости.
	InitialPayment float64         `json:"initial_payment"` // Первоначальный взнос.
	Months         int             `json:"months"`          // Срок кредита в месяцах.
	Program        MortgageProgram `json:"program"`         // Выбранная ипотечная программа.
}

// MortgageResponse представляет успешный JSON-ответ сервиса,
// включающий исходные параметры запроса, программу и агрегаты.
type MortgageResponse struct {
	Aggregates MortgageAggregates `json:"aggregates"` // Рассчитанные агрегаты кредита.
	Params     MortgageParams     `json:"params"`     // Входные параметры кредита.
	Program    MortgageProgram    `json:"program"`    // Выбранная ипотечная программа.
}

// DateFormat задаёт формат даты для сериализации LastPaymentDate.
const DateFormat = "2006-01-02"
