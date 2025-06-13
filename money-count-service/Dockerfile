FROM golang:1.24-alpine AS builder

# устанавливаем пакеты
RUN apk add --no-cache git protobuf protobuf-dev

# устанавливаем protoc
RUN go install google.golang.org/protobuf/cmd/protoc-gen-go@latest && \
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# рабочая директория
WORKDIR /app

# копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# загружаем зависимости
RUN go mod download

# копируем gen_backend_proto.sh explicitly
COPY gen_backend_proto.sh .
RUN chmod +x gen_backend_proto.sh
RUN ls -la /app  # Debug: Confirm file presence

# копируем остальной исходный код
COPY . .

# запускаем скрипт генерации
RUN ./gen_backend_proto.sh

# собираем приложение
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

# Экспонируем порты
EXPOSE 50051 8080

# Healthcheck
HEALTHCHECK --interval=30s --timeout=5s --start-period=5s --retries=3 \
    CMD timeout 5s sh -c "</dev/tcp/localhost/50051" || exit 1

# Запускаем приложение
CMD ["./money-service"]