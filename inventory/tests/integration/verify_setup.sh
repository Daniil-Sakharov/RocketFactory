#!/bin/bash
# Скрипт для проверки готовности окружения к запуску интеграционных тестов

set -e

echo "🔍 Проверка готовности к запуску интеграционных тестов..."
echo ""

# Цвета для вывода
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Счетчики
PASSED=0
FAILED=0

# Функция для проверки
check() {
    local name=$1
    local command=$2
    
    echo -n "Проверка: $name... "
    if eval "$command" > /dev/null 2>&1; then
        echo -e "${GREEN}✅ OK${NC}"
        ((PASSED++))
        return 0
    else
        echo -e "${RED}❌ FAILED${NC}"
        ((FAILED++))
        return 1
    fi
}

# Функция для вывода предупреждения
warn() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

# 1. Проверка Docker
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "1️⃣  Docker"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if check "Docker установлен" "command -v docker"; then
    if check "Docker daemon запущен" "docker ps"; then
        DOCKER_VERSION=$(docker --version)
        echo "   Docker version: $DOCKER_VERSION"
    else
        warn "Docker установлен, но daemon не запущен. Запустите: sudo systemctl start docker"
    fi
else
    warn "Docker не установлен. Установите Docker для запуска интеграционных тестов."
    warn "Инструкции: https://docs.docker.com/get-docker/"
fi
echo ""

# 2. Проверка Go
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "2️⃣  Go"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if check "Go установлен" "command -v go"; then
    GO_VERSION=$(go version)
    echo "   Go version: $GO_VERSION"
else
    warn "Go не установлен"
fi
echo ""

# 3. Проверка .env файлов
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "3️⃣  Environment Files"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"

# Определяем корневую директорию проекта
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "$SCRIPT_DIR/../../.." && pwd )"

BASE_ENV_FILE="$PROJECT_ROOT/deploy/env/.env"
INVENTORY_ENV_FILE="$PROJECT_ROOT/deploy/compose/inventory/.env"

if check "Базовый .env файл существует" "test -f $BASE_ENV_FILE"; then
    echo "   Путь: $BASE_ENV_FILE"
else
    warn "Создайте базовый .env файл: cp deploy/env/.env.template deploy/env/.env"
fi

if check "Inventory .env файл существует" "test -f $INVENTORY_ENV_FILE"; then
    echo "   Путь: $INVENTORY_ENV_FILE"
    
    # Проверяем наличие важных переменных
    if grep -q "MONGO_DATABASE=inventory" "$INVENTORY_ENV_FILE"; then
        echo -e "   ${GREEN}✓${NC} MONGO_DATABASE настроена"
    else
        warn "Переменная MONGO_DATABASE не найдена в .env файле"
    fi
    
    if grep -q "GRPC_PORT=50051" "$INVENTORY_ENV_FILE"; then
        echo -e "   ${GREEN}✓${NC} GRPC_PORT настроен"
    else
        warn "Переменная GRPC_PORT не найдена в .env файле"
    fi
else
    warn "Создайте .env файл для inventory сервиса:"
    warn "  cd $PROJECT_ROOT && task env:generate"
    warn "  или вручную: cd deploy/env && SERVICES=inventory ./generate-env.sh"
fi
echo ""

# 4. Проверка Dockerfile
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "4️⃣  Dockerfile"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
DOCKERFILE="$PROJECT_ROOT/deploy/docker/inventory/Dockerfile"
if check "Dockerfile существует" "test -f $DOCKERFILE"; then
    echo "   Путь: $DOCKERFILE"
else
    warn "Dockerfile не найден: $DOCKERFILE"
fi
echo ""

# 5. Проверка зависимостей Go
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "5️⃣  Go Dependencies"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
if check "testcontainers-go установлен" "go list -m github.com/testcontainers/testcontainers-go"; then
    VERSION=$(go list -m github.com/testcontainers/testcontainers-go | awk '{print $2}')
    echo "   Version: $VERSION"
fi

if check "ginkgo установлен" "go list -m github.com/onsi/ginkgo/v2"; then
    VERSION=$(go list -m github.com/onsi/ginkgo/v2 | awk '{print $2}')
    echo "   Version: $VERSION"
fi
echo ""

# Итоговый результат
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo "📊 Результаты"
echo "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━"
echo -e "${GREEN}✅ Passed: $PASSED${NC}"
echo -e "${RED}❌ Failed: $FAILED${NC}"
echo ""

if [ $FAILED -eq 0 ]; then
    echo -e "${GREEN}🎉 Все проверки пройдены! Можно запускать тесты:${NC}"
    echo ""
    echo "  # Через Taskfile:"
    echo "  task test-integration MODULES=inventory"
    echo ""
    echo "  # Или напрямую:"
    echo "  go test -v -timeout=20m -tags=integration ./inventory/tests/integration/..."
    echo ""
    exit 0
else
    echo -e "${YELLOW}⚠️  Некоторые проверки не пройдены. Исправьте проблемы выше.${NC}"
    echo ""
    echo "💡 Быстрая настройка:"
    echo "  1. Установите Docker: https://docs.docker.com/get-docker/"
    echo "  2. Запустите: task env:generate"
    echo "  3. Запустите этот скрипт снова для проверки"
    echo ""
    exit 1
fi
