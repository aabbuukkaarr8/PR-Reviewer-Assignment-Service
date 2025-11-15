.PHONY: help build run test clean docker-build docker-up docker-down docker-logs fmt deps migrate-up migrate-down

BINARY_NAME=apiserver
BINARY_PATH=./cmd/apiserver
CONFIG_PATH=configs/apiserver.toml
DOCKER_COMPOSE=docker-compose
GO=go

GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m

help: ## Показать справку по командам
	@echo "$(GREEN)Доступные команды:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

deps: ## Установить зависимости
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	$(GO) mod download
	$(GO) mod tidy

build: ## Собрать бинарник
	@echo "$(GREEN)Сборка бинарника...$(NC)"
	$(GO) build -o $(BINARY_NAME) $(BINARY_PATH)
	@echo "$(GREEN)Бинарник создан: $(BINARY_NAME)$(NC)"

run: ## Запустить сервис локально
	@echo "$(GREEN)Запуск сервиса...$(NC)"
	$(GO) run $(BINARY_PATH) -config-path=$(CONFIG_PATH)

test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	$(GO) test -v ./...

test-coverage: ## Запустить тесты с покрытием
	@echo "$(GREEN)Запуск тестов с покрытием...$(NC)"
	$(GO) test -v -coverprofile=coverage.out ./...
	$(GO) tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)Отчет о покрытии создан: coverage.html$(NC)"

fmt: ## Форматировать код
	@echo "$(GREEN)Форматирование кода...$(NC)"
	$(GO) fmt ./...

clean: ## Очистить сгенерированные файлы
	@echo "$(GREEN)Очистка...$(NC)"
	rm -f $(BINARY_NAME)
	rm -f coverage.out coverage.html
	$(GO) clean

docker-build: ## Собрать Docker образы
	@echo "$(GREEN)Сборка Docker образов...$(NC)"
	$(DOCKER_COMPOSE) build

docker-up: ## Запустить сервисы через Docker Compose
	@echo "$(GREEN)Запуск сервисов через Docker Compose...$(NC)"
	$(DOCKER_COMPOSE) up -d
	@echo "$(GREEN)Сервисы запущены. API доступен на http://localhost:8080$(NC)"

docker-down: ## Остановить сервисы Docker Compose
	@echo "$(GREEN)Остановка сервисов...$(NC)"
	$(DOCKER_COMPOSE) down

docker-down-volumes: ## Остановить сервисы и удалить volumes
	@echo "$(GREEN)Остановка сервисов и удаление volumes...$(NC)"
	$(DOCKER_COMPOSE) down -v

docker-logs: ## Показать логи Docker Compose
	$(DOCKER_COMPOSE) logs -f

docker-restart: ## Перезапустить сервисы
	@echo "$(GREEN)Перезапуск сервисов...$(NC)"
	$(DOCKER_COMPOSE) restart

migrate-up: ## Применить миграции БД
	@echo "$(GREEN)Миграции применяются автоматически при запуске БД$(NC)"

migrate-down: ## Откатить миграции БД
	@echo "$(YELLOW)Для отката миграций используйте инструмент migrate или psql$(NC)"

dev: docker-up ## Запустить окружение для разработки
	@echo "$(GREEN)Окружение для разработки готово!$(NC)"

stop: docker-down ## Остановить окружение

all: clean deps build test ## Выполнить полную сборку
	@echo "$(GREEN)Полная сборка завершена!$(NC)"

