# Многоэтапная сборка для оптимизации размера образа
FROM golang:1.21-alpine AS builder

# Устанавливаем необходимые пакеты для сборки
RUN apk add --no-cache git protobuf protobuf-dev

# Устанавливаем protoc плагины для Go
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Генерируем proto файлы
RUN chmod +x gen_backend_proto.sh && ./gen_backend_proto.sh

# Собираем приложение
RUN CGO_ENABLED=0 go build -ldflags="-w -s" -o bin/money-service cmd/server/main.go

# Финальный образ
FROM alpine:latest

# Устанавливаем необходимые пакеты
RUN apk add --no-cache ca-certificates tzdata

# Создаем пользователя для безопасности
RUN adduser -D -s /bin/sh appuser

# Создаем рабочую директорию
WORKDIR /app

# Копируем бинарник из builder
COPY --from=builder /app/bin/money-service .

# Устанавливаем права на исполнение
RUN chmod +x money-service

# Переключаемся на пользователя appuser
USER appuser

# Экспонируем порт (измените если нужен другой)
EXPOSE 50051

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD timeout 5s sh -c "</dev/tcp/localhost/50051" || exit 1

# Запускаем приложение
CMD ["./money-service"]