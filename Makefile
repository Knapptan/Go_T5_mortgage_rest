# Makefile

# Переменные
NAME=mortgage-calculator
IMAGE_NAME=mortgage-app
COVER_PROFILE=coverage.out
TEST_PACKAGES=$$(go list ./... \
  | grep -v "/cmd" \
  | grep -v "/internal/models")

# Цели по умолчанию
.DEFAULT_GOAL := help

.PHONY: help build run lint test test-coverage test-race docker-build docker-run docker-stop clean

help:
	@grep -E '^[a-zA-Z_-]+:.*?##' Makefile | \
      awk 'BEGIN {FS = ":.*?##"} {printf "  %-15s %s\n", $$1, $$2}'

build: ## Собрать проект
	go build -o bin/$(NAME) ./cmd/

run: ## Запустить сервер
	go run ./cmd/

lint: ## Запустить линтер по всему проекту
	golangci-lint run

test: ## Запустить все тесты
	go test -v $(TEST_PACKAGES)

test-coverage: ## Покрытие >80% для бизнес-пакетов
	@go test $(TEST_PACKAGES) -coverprofile=$(COVER_PROFILE)
	@go tool cover -func=$(COVER_PROFILE)
	@rm -f $(COVER_PROFILE)

test-race: ## Проверить на гонки данных
	go test -race $(TEST_PACKAGES)

docker-build: ## Сборка Docker-образа
	docker build -t $(NAME) .

docker-run: ## Запуск контейнера в фоне
	docker run -d -p 8080:8080 --name $(IMAGE_NAME) $(NAME)

docker-stop: ## Остановка и удаление контейнера
	docker stop $(IMAGE_NAME) && docker rm $(IMAGE_NAME)

clean: ## Очистить артефакты сборки
	@rm -rf bin/ $(COVER_PROFILE)