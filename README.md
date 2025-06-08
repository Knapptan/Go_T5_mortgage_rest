# Сервис расчета параметров ипотеки (Ипотечный калькулятор)

Микросервис для расчета параметров ипотеки с кэшированием результатов.

## 🚀 Быстрый старт

### Предварительные требования

- Установленный Docker
- Установленный Make (опционально)

### Запуск с помощью Docker

```bash
# Сборка Docker-образа
make docker-build

# Запуск контейнера
make docker-run

# Остановка контейнера
make docker-stop
```

### Прямой запуск (без Docker)

```bash
# Установите зависимости
go mod download

# Запустите сервис
go run ./cmd
```

Сервис будет доступен по адресу: http://localhost:8080

## 📡 API Endpoints

### 1. Расчет параметров ипотеки

POST /execute

Пример запроса:

```json
{
	"object_cost": 5000000,
	"initial_payment": 1000000,
	"months": 240,
	"program": {
		"salary": true
	}
}
```

Пример успешного ответа (200 OK):

```json
{
	"result": {
		"params": {
			"object_cost": 5000000,
			"initial_payment": 1000000,
			"months": 240
		},
		"program": {
			"salary": true,
			"military": false,
			"base": false
		},
		"aggregates": {
			"rate": 8,
			"loan_sum": 4000000,
			"monthly_payment": 33457.6,
			"overpayment": 4029824,
			"last_payment_date": "2045-06-08"
		}
	}
}
```

Возможные ошибки (400 Bad Request):

```json
{
	"error": "choose program"
}
```

```json
{
	"error": "choose only 1 program"
}
```

```json
{
	"error": "the initial payment should be more"
}
```

```json
{
	"error": "object cost and initial payment must be positive"
}
```

```json
{
	"error": "loan duration must be at least 1 month"
}
```

```json
{
	"error": "initial payment cannot exceed object cost"
}
```

### 2. Получение кэшированных расчетов

GET /cache

Пример успешного ответа (200 OK):

```json
[
	{
		"id": 1,
		"params": {
			"object_cost": 5000000,
			"initial_payment": 1000000,
			"months": 240
		},
		"program": {
			"salary": true,
			"military": false,
			"base": false
		},
		"aggregates": {
			"rate": 8,
			"loan_sum": 4000000,
			"monthly_payment": 33457.6,
			"overpayment": 4029824,
			"last_payment_date": "2045-06-08"
		}
	}
]
```

Ошибка при пустом кэше (400 Bad Request):

```json
{
	"error": "empty cache"
}
```

## 🛠 Техническое задание

### Функциональность

#### Расчет параметров ипотеки:

- Процентная ставка по программе кредитования

- Сумма кредита

- Ежемесячный аннуитетный платеж

- Переплата за весь срок

- Дата последнего платежа

- Поддержка 3 программ кредитования:

- Корпоративные клиенты (8%)

- Военная ипотека (9%)

- Базовая программа (10%)

- Проверка первоначального взноса (не менее 20%)

- Кэширование результатов расчетов

- Логирование запросов

### Технические требования

- Web-framework: Gin

- Конфигурация: через config.yml (порт 8080)

- Тестирование:

- Покрытие >80%

- make test для запуска тестов

#### Линтинг:

- golangci-lint с кастомной конфигурацией

- make lint для проверки

#### Вендоринг зависимостей: go mod vendor

#### Docker:

- Образ на базе Alpine Linux

- Размер <30MB (фактический: ~22.5MB)

#### Makefile с командами для сборки, тестирования и запуска

## 🧩 Программы кредитования

- Программа для корпоративных клиентов: 8%

- Военная ипотека: 9%

- Базовая программа: 10%

## ⚙️ Системные команды

Полный список команд для работы с проектом:

```bash
# Показать справку по всем командам
make help

# Собрать бинарный файл проекта
make build

# Запустить сервер в режиме разработки
make run

# Запустить линтер для проверки кода
make lint

# Запустить все тесты
make test

# Проверить покрытие кода тестами
make test-coverage

# Проверить наличие гонок данных (race conditions)
make test-race

# Собрать Docker-образ
make docker-build

# Запустить Docker-контейнер
make docker-run

# Остановить и удалить Docker-контейнер
make docker-stop

# Очистить артефакты сборки
make clean
```



## 📂 Структура проекта

```tree
.
├── Dockerfile                # Конфигурация Docker
├── Makefile                  # Автоматизация задач
├── README.md                 # Документация
├── cmd
│   └── main.go               # Точка входа
├── config
│   ├── config.go             # Загрузка конфигурации
│   └── config_test.go        # Тесты конфигурации
├── config.yml                # Файл конфигурации
├── go.mod
├── go.sum
├── internal
│   ├── cache                 # Реализация кэша
│   │   ├── cache.go
│   │   └── cache_test.go
│   ├── calculator            # Логика расчетов
│   │   ├── calculator.go
│   │   └── calculator_test.go
│   ├── http                  # HTTP обработчики
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── middleware.go     # Логирование запросов
│   │   └── routes.go
│   └── models                # Модели данных
│       └── models.go
└── vendor                    # Вендоринг зависимостей
````

## 📈 Технические показатели

Покрытие тестами: 94.5%

Размер Docker-образа: 22.5MB

Время обработки запроса: <5ms

Потокобезопасность: Поддержка конкурентных запросов

Грейсфул шатдаун: Корректное завершение работы
