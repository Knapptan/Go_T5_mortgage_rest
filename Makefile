# Makefile для управления проектом

# Переменные
COVER_PROFILE=coverage.out
TEST_PACKAGES=
NAME=mortgage-calculator
IMAGE_NAME=mortgage-app

# Цели по умолчанию
.DEFAULT_GOAL := help

.PHONY: help
help: ## Показать справку по целям
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Основные цели
.PHONY: build
build: ## Собрать проект
	go build -o bin/$(NAME) ./cmd/

.PHONY: run
run: ## Запустить сервер
	go run ./cmd/

.PHONY: lint
lint: ## Запустить линтер
	golangci-lint run -c .golangci.yml ./cmd/

.PHONY: docker-build
docker-build: ## Сборка Docker-образа
	docker build -t $(NAME) .

.PHONY: docker-run
docker-run: ## Запуск контейнера в фоне, проброс порта
docker run -d -p 8080:8080 --name $(IMAGE_NAME) $(NAME)

.PHONY: docker-stop
docker-stop: ## Остановка и удаление контейнера
	docker stop $(IMAGE_NAME)
	docker rm $(IMAGE_NAME)

.PHONY: test
test: ## Запустить все тесты
	go test -v $(TEST_PACKAGES)

.PHONY: test-coverage
test-coverage: ## Запуск тестов для проверки покрытия
	go test ./... -coverprofile=coverage.out
  go tool cover -func=coverage.out
  go tool cover -html=coverage.out -o coverage.html

.PHONY: test-race
test-race: ## Проверить на гонки данных
	go test -race $(TEST_PACKAGES)

.PHONY: clean
clean: ## Очистить артефакты сборки
	rm -rf bin/
	rm -f $(COVER_PROFILE)