#!/bin/sh
set -e

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}=== Генерация gRPC-Web файлов (чистый фронтенд) ===${NC}"

# 1. Проверка proto-файла
PROTO_FILE="proto/money_service.proto"
if [ ! -f "$PROTO_FILE" ]; then
    echo -e "${RED}ОШИБКА: Файл $PROTO_FILE не найден в папке фронтенда!${NC}"
    echo -e "Полный путь: $(pwd)/$PROTO_FILE"
    echo -e "Убедитесь, что файл скопирован в правильную директорию"
    exit 1
fi

# 2. Подготовка директории
GEN_DIR="gen"
echo -e "Подготовка директории $GEN_DIR/"
mkdir -p "$GEN_DIR"
rm -f "$GEN_DIR"/*.js 2>/dev/null || true

# 3. Генерация файлов
echo -e "Генерация из локального proto-файла:"
protoc \
  -I=proto \
  --grpc-web_out=import_style=commonjs,mode=grpcwebtext:"$GEN_DIR" \
  "$PROTO_FILE"

# 4. Проверка результата
if [ -z "$(ls -A $GEN_DIR/*.js 2>/dev/null)" ]; then
    echo -e "${RED}ОШИБКА: Файлы не были сгенерированы!${NC}"
    exit 1
fi

echo -e "${GREEN}Успешно сгенерированы:${NC}"
ls -1 "$GEN_DIR"/*.js | xargs -n1 basename

echo -e "${GREEN}=== Генерация завершена ===${NC}"