// Package models содержит модели данных
package models

// Структура запроса на сервис
type MortgageRequest struct {
	ObjectCost     float64  `json:"object_cost"`     // стоимость объекта
	InitialPayment float64  `json:"initial_payment"` // первоначальный взнос
	Months         int      `json:"months"`          // срок
	Program        struct { // блок программы кредита
		Salary   bool `json:"salary"`   // программа для корпоративных клиентов
		Military bool `json:"military"` // военная ипотека
		Base     bool `json:"base"`     // базовая программа
	} `json:"program"`
}

// Структура успешного ответа сервиса
type MortgageResponse struct {
	Params struct { // запрашиваемые параметры кредита
		ObjectCost     float64 `json:"object_cost"`     // стоимость объекта
		InitialPayment float64 `json:"initial_payment"` // первоначальный взнос
		Months         int     `json:"months"`          // срок
	} `json:"params"`
	Program struct { // блок программы кредита
		Salary   bool `json:"salary"`   // программа для корпоративных клиентов
		Military bool `json:"military"` // военная ипотека
		Base     bool `json:"base"`     // базовая программа
	} `json:"program"`
	Aggregates struct { // блок с агрегатами
		Rate            float64 `json:"rate"`              // годовая процентная ставка
		LoanSum         float64 `json:"loan_sum"`          // сумма кредита
		MonthlyPayment  float64 `json:"monthly_payment"`   // аннуитетный ежемесячный платеж
		Overpayment     float64 `json:"overpayment"`       // переплата за весь срок кредита
		LastPaymentDate string  `json:"last_payment_date"` // последняя дата платежа
	} `json:"aggregates"`
}

// Структура для кэширования запросов
type CacheItem struct {
	ID       int              // id расчета в кэше
	Request  MortgageRequest  // Запрос на сервис
	Response MortgageResponse // Ответ сервиса
}
