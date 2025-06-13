# Переменные
SERVICE_NAME = money-count-service
DOCKER_IMAGE = wb-$(SERVICE_NAME)
DOCKER_TAG = latest
DOCKER_NETWORK = clickhouse-network
GRPC_PORT = 50051
WEB_PORT = 8080
CONTAINER_NAME = wb-$(SERVICE_NAME)

# Go переменные
GO_BUILD_ENV = CGO_ENABLED=0 GOOS=linux
GO_FILES = $(shell find . -type f -name '*.go')
PROTO_FILES = proto/money_service.proto

# Цвета для вывода
GREEN = \033[0;32m
YELLOW = \033[1;33m
RED = \033[0;31m
NC = \033[0m # No Color

.PHONY: help proto build test docker-build docker-run docker-stop docker-logs clean deps lint

# По умолчанию показываем help
help: ## Показать справку
	@echo "$(GREEN)Доступные команды:$(NC)"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Установка зависимостей
deps: ## Установить зависимости Go
	@echo "$(GREEN)Установка зависимостей...$(NC)"
	@go mod download
	@go mod tidy
	@echo "$(GREEN)Зависимости установлены$(NC)"

# Генерация proto файлов
proto: ## Генерировать Go файлы из proto
	@echo "$(GREEN)Генерация proto файлов...$(NC)"
	@chmod +x gen_backend_proto.sh
	@./gen_backend_proto.sh
	@echo "$(GREEN)Proto файлы сгенерированы$(NC)"

# Линтинг кода
lint: ## Запустить линтер
	@echo "$(GREEN)Запуск линтера...$(NC)"
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run; \
	else \
		echo "$(YELLOW)golangci-lint не установлен. Используем go vet...$(NC)"; \
		go vet ./...; \
	fi

# Тестирование
test: ## Запустить тесты
	@echo "$(GREEN)Запуск тестов...$(NC)"
	@go test -v ./...

# Локальная сборка
build: proto ## Собрать приложение локально
	@echo "$(GREEN)Сборка приложения...$(NC)"
	@mkdir -p bin
	@$(GO_BUILD_ENV) go build -ldflags="-w -s" -o bin/$(SERVICE_NAME) cmd/server/main.go
	@echo "$(GREEN)Приложение собрано: bin/$(SERVICE_NAME)$(NC)"

# Запуск локально
run: build ## Запустить приложение локально
	@echo "$(GREEN)Запуск приложения...$(NC)"
	@./bin/$(SERVICE_NAME)

# Docker команды
docker-build: ## Собрать Docker образ
	@echo "$(GREEN)Сборка Docker образа...$(NC)"
	@docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .
	@echo "$(GREEN)Docker образ собран: $(DOCKER_IMAGE):$(DOCKER_TAG)$(NC)"

docker-run: docker-stop ## Запустить контейнер
	@echo "$(GREEN)Запуск контейнера...$(NC)"
	@docker run -d \
		--name $(CONTAINER_NAME) \
		--network $(DOCKER_NETWORK) \
		-p $(GRPC_PORT):$(GRPC_PORT) \
		-p $(WEB_PORT):$(WEB_PORT) \
		--restart unless-stopped \
		$(DOCKER_IMAGE):$(DOCKER_TAG)
	@echo "$(GREEN)Контейнер запущен: $(CONTAINER_NAME)$(NC)"
	@echo "$(YELLOW)gRPC доступен на порту $(GRPC_PORT)$(NC)"
	@echo "$(YELLOW)gRPC-Web доступен на порту $(WEB_PORT)$(NC)"

docker-stop: ## Остановить и удалить контейнер
	@echo "$(GREEN)Остановка контейнера...$(NC)"
	@-docker stop $(CONTAINER_NAME) 2>/dev/null || true
	@-docker rm $(CONTAINER_NAME) 2>/dev/null || true

docker-logs: ## Показать логи контейнера
	@docker logs -f $(CONTAINER_NAME)

docker-shell: ## Войти в контейнер
	@docker exec -it $(CONTAINER_NAME) sh

# Комплексные команды
up: docker-build docker-run ## Полная пересборка и запуск
	@echo "$(GREEN)Сервис запущен!$(NC)"

restart: docker-stop docker-run ## Перезапустить контейнер
	@echo "$(GREEN)Сервис перезапущен!$(NC)"

# Очистка
clean: ## Очистить сгенерированные файлы и образы
	@echo "$(GREEN)Очистка...$(NC)"
	@rm -rf bin/
	@rm -f pkg/proto/*.pb.go
	@-docker rmi $(DOCKER_IMAGE):$(DOCKER_TAG) 2>/dev/null || true
	@echo "$(GREEN)Очистка завершена$(NC)"

# Разработка
dev-setup: deps proto ## Настройка среды разработки
	@echo "$(GREEN)Среда разработки настроена$(NC)"

# Статус
status: ## Показать статус контейнера
	@echo "$(GREEN)Статус контейнера:$(NC)"
	@docker ps --filter name=$(CONTAINER_NAME) --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}"

# Полная пересборка для продакшена
production: clean deps lint test docker-build ## Полная сборка для продакшена
	@echo "$(GREEN)Продакшен сборка завершена!$(NC)"

# Docker Compose интеграция
compose-up: ## Запустить все сервисы через docker-compose
	@echo "$(GREEN)Запуск всех сервисов...$(NC)"
	@docker-compose up -d

compose-down: ## Остановить все сервисы
	@echo "$(GREEN)Остановка всех сервисов...$(NC)"
	@docker-compose down

compose-logs: ## Показать логи всех сервисов
	@docker-compose logs -f