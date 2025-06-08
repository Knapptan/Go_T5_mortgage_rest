# Ипотечный калькулятор

## Запуск

## Функциональность:

- Расчет всех параметров ипотеки

- Три программы кредитования с разными ставками

- Проверка первоначального взноса (20%)

- Кэширование результатов

- Эндпоинты /execute (POST) и /cache (GET)

- Обработка ошибок

- Middleware для логирования

## Технические требования:

- Web-framework (Gin)

- Конфигурация через config.yml

- Покрытие тестами >80%

- Проходит линтер golangci-lint

- Вендоринг зависимостей (go mod vendor)

- Dockerfile

- Вес образа <30MB (Alpine-based образ ~12MB)

## Makefile с командами

### Makefile команды:

- make test - запуск тестов

- make lint - запуск линтера

- make docker-build - сборка образа

- make docker-run - запуск контейнера

- make docker-stop - остановка контейнера

Структура проекта

```tree
.
├── Dockerfile
├── Makefile
├── README.md
├── cmd
│   └── main.go
├── config
│   ├── config.go
│   └── config_test.go
├── config.yml
├── go.mod
├── go.sum
├── internal
│   ├── cache
│   │   ├── cache.go
│   │   └── cache_test.go
│   ├── calculator
│   │   ├── calculator.go
│   │   └── calculator_test.go
│   ├── http
│   │   ├── handler.go
│   │   ├── handler_test.go
│   │   ├── http_test.go
│   │   ├── middleware.go
│   │   └── routes.go
│   └── models
│       └── models.go
└── vendor
    └── ... # вендоринг зависимостей
```
