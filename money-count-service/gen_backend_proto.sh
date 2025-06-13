#!/bin/bash

# Скрипт для генерации Go файлов из proto (только для бэкенда)

set -e  # Прерывать выполнение при любой ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE} Генерация proto файлов для Go бэкенда${NC}"

# Проверяем наличие protoc
if ! command -v protoc &> /dev/null; then
    echo -e "${RED} protoc не найден! Установите Protocol Buffers compiler${NC}"
    echo -e "${YELLOW}macOS: brew install protobuf${NC}"
    echo -e "${YELLOW}Ubuntu: sudo apt install protobuf-compiler${NC}"
    echo -e "${YELLOW}Windows: https://grpc.io/docs/protoc-installation/${NC}"
    exit 1
fi

# Проверяем наличие protoc-gen-go
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${YELLOW} protoc-gen-go не найден! Устанавливаю...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
    echo -e "${GREEN}protoc-gen-go установлен${NC}"
fi

# Проверяем наличие protoc-gen-go-grpc
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${YELLOW}protoc-gen-go-grpc не найден! Устанавливаю...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
    echo -e "${GREEN}protoc-gen-go-grpc установлен${NC}"
fi

# Проверяем существование proto файла
if [ ! -f "proto/money_service.proto" ]; then
    echo -e "${RED}Файл proto/money_service.proto не найден!${NC}"
    exit 1
fi

# Создаём директорию если её нет
mkdir -p pkg/proto

echo -e "${YELLOW}Очищаю старые Go файлы...${NC}"
rm -f pkg/proto/*.pb.go

echo -e "${YELLOW} Генерирую Go файлы из proto...${NC}"

# Генерируем Go код
protoc \
    -I proto \
    --go_out=pkg/proto --go_opt=paths=source_relative \
    --go-grpc_out=pkg/proto --go-grpc_opt=paths=source_relative \
    proto/money_service.proto

# Проверяем результат
if [ -f "pkg/proto/money_service.pb.go" ] && [ -f "pkg/proto/money_service_grpc.pb.go" ]; then
    echo -e "${GREEN}Генерация успешно завершена!${NC}"
    echo -e "${BLUE}Сгенерированные файлы:${NC}"
    ls -la pkg/proto/*.pb.go

    # Показываем размер файлов
    echo -e "${BLUE}Размеры файлов:${NC}"
    du -h pkg/proto/*.pb.go

    echo -e "${GREEN}Готово! Теперь можно компилировать Go сервер${NC}"
else
    echo -e "${RED}Ошибка генерации! Файлы не созданы${NC}"
    exit 1
fi

# Дополнительная информация
echo -e "${BLUE} Совет: Для компиляции сервера используйте:${NC}"
echo -e "   ${YELLOW}go run cmd/server/main.go${NC}"
echo -e "${BLUE} Для сборки бинарника:${NC}"
echo -e "   ${YELLOW}go build -o bin/money-service cmd/server/main.go${NC}"