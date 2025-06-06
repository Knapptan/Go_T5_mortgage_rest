# Makefile для управления проектом

# Переменные
GO=go
COVER_PROFILE=coverage.out
# TEST_PACKAGES=

# Цели по умолчанию
.DEFAULT_GOAL := help

.PHONY: help
help: ## Показать справку по целям
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Основные цели
.PHONY: build
build: ## Собрать проект
	$(GO) build -o bin/server ./cmd/

.PHONY: run
run: ## Запустить сервер
	$(GO) run ./cmd/

# Цели тестирования
.PHONY: test
test: ## Запустить все тесты
	$(GO) test -v $(TEST_PACKAGES)

.PHONY: test-race
test-race: ## Проверить на гонки данных
	$(GO) test -race $(TEST_PACKAGES)

.PHONY: clean
clean: ## Очистить артефакты сборки
	rm -rf bin/
	rm -f $(COVER_PROFILE)