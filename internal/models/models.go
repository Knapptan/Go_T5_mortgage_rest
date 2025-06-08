// Package models содержит модели данных
package models

// Структура запрашиваемых параметров кредита
type MortgageParams struct {
	ObjectCost     float64 `json:"object_cost"`     // стоимость объекта
	InitialPayment float64 `json:"initial_payment"` // первоначальный взнос
	Months         int     `json:"months"`          // срок
}

// Структура программы кредита
type MortgageProgram struct {
	Salary   bool `json:"salary"`   // программа для корпоративных клиентов
	Military bool `json:"military"` // военная ипотека
	Base     bool `json:"base"`     // базовая программа
}

// Структура с агрегатами
type MortgageAggregates struct {
	Rate            float64 `json:"rate"`              // годовая процентная ставка
	LoanSum         float64 `json:"loan_sum"`          // сумма кредита
	MonthlyPayment  float64 `json:"monthly_payment"`   // аннуитетный ежемесячный платеж
	Overpayment     float64 `json:"overpayment"`       // переплата за весь срок кредита
	LastPaymentDate string  `json:"last_payment_date"` // последняя дата платежа
}

// Структура для кэширования запросов
type MortgageInfoResponse struct {
	ID         int                // id расчета в кэше
	Params     MortgageParams     `json:"params"`     // запрашиваемые параметры кредита
	Program    MortgageProgram    `json:"program"`    // блок программы кредита
	Aggregates MortgageAggregates `json:"aggregates"` // блок с агрегатами
}

// Структура запроса на сервис
type MortgageRequest struct {
	ObjectCost     float64         `json:"object_cost"`     // стоимость объекта
	InitialPayment float64         `json:"initial_payment"` // первоначальный взнос
	Months         int             `json:"months"`          // срок
	Program        MortgageProgram `json:"program"`         // блок программы кредита
}

// Структура успешного ответа сервиса
type MortgageResponse struct {
	Params     MortgageParams     `json:"params"`     // запрашиваемые параметры кредита
	Program    MortgageProgram    `json:"program"`    // блок программы кредита
	Aggregates MortgageAggregates `json:"aggregates"` // блок с агрегатами
}

// Константа временного формата
const DateFormat = "2006-01-02"
