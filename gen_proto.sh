#!/bin/bash

# Скрипт для генерации Go и JS файлов из proto

set -e  # Прерывать выполнение при любой ошибке

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🚀 Генерация proto файлов...${NC}"

# Проверяем наличие protoc
if ! command -v protoc &> /dev/null; then
    echo -e "${RED}❌ protoc не найден! Установите Protocol Buffers compiler${NC}"
    echo -e "${YELLOW}macOS: brew install protobuf${NC}"
    echo -e "${YELLOW}Ubuntu: sudo apt install protobuf-compiler${NC}"
    exit 1
fi

# Проверяем наличие protoc-gen-go
if ! command -v protoc-gen-go &> /dev/null; then
    echo -e "${RED}❌ protoc-gen-go не найден! Устанавливаю...${NC}"
    go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
fi

# Проверяем наличие protoc-gen-go-grpc
if ! command -v protoc-gen-go-grpc &> /dev/null; then
    echo -e "${RED}❌ protoc-gen-go-grpc не найден! Устанавливаю...${NC}"
    go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
fi

# Создаём директории если их нет
mkdir -p pkg/proto
mkdir -p web/js/gen

echo -e "${YELLOW}📁 Очищаю старые файлы...${NC}"
rm -f pkg/proto/*.pb.go
rm -f web/js/gen/*

echo -e "${YELLOW}🔧 Генерирую Go файлы...${NC}"
protoc \
    --go_out=. \
    --go_opt=paths=source_relative \
    --go-grpc_out=. \
    --go-grpc_opt=paths=source_relative \
    proto/money_service.proto

echo -e "${YELLOW}🌐 Генерирую JavaScript файлы для gRPC-Web...${NC}"
protoc \
    --js_out=import_style=commonjs:web/js/gen \
    --grpc-web_out=import_style=commonjs,mode=grpcwebtext:web/js/gen \
    proto/money_service.proto

echo -e "${GREEN}✅ Генерация завершена!${NC}"
echo -e "${BLUE}📄 Сгенерированные файлы:${NC}"
echo -e "  ${GREEN}Go:${NC}"
ls -la pkg/proto/*.pb.go 2>/dev/null || echo -e "    ${RED}Нет файлов${NC}"
echo -e "  ${GREEN}JavaScript:${NC}"
ls -la web/js/gen/* 2>/dev/null || echo -e "    ${RED}Нет файлов${NC}"

echo -e "${BLUE}🎉 Готово! Теперь можно использовать обновлённые proto файлы${NC}"