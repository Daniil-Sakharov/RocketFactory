#!/bin/bash

# 🧪 Скрипт для запуска тестов Inventory Service

set -e

echo "🚀 Запуск тестов для Inventory Service..."

# Проверяем, что мы в правильной директории
if [ ! -f "go.mod" ]; then
    echo "❌ Ошибка: go.mod не найден. Запустите скрипт из корня проекта."
    exit 1
fi

echo ""
echo "📋 Доступные типы тестов:"
echo "1. Unit тесты (рекомендуется для CI/CD)"
echo "2. Интеграционные тесты (требуют Docker)"
echo "3. Все тесты"

read -p "Выберите тип тестов (1-3): " choice

case $choice in
    1)
        echo ""
        echo "🧪 Запуск unit тестов..."
        cd inventory/tests/unit
        go test -v -timeout=20m -tags=unit .
        ;;
    2)
        echo ""
        echo "🧪 Запуск интеграционных тестов..."
        echo "⚠️  Внимание: Требуется Docker!"
        
        # Проверяем доступность Docker
        if ! command -v docker &> /dev/null; then
            echo "❌ Docker не установлен!"
            exit 1
        fi
        
        if ! docker ps &> /dev/null; then
            echo "❌ Docker не запущен или недоступен!"
            echo "💡 Попробуйте: sudo systemctl start docker"
            exit 1
        fi
        
        cd inventory/tests/integration
        go test -v -timeout=20m -tags=integration .
        ;;
    3)
        echo ""
        echo "🧪 Запуск всех тестов..."
        
        echo "▶️  Unit тесты..."
        cd inventory/tests/unit
        go test -v -timeout=20m -tags=unit .
        
        echo ""
        echo "▶️  Интеграционные тесты..."
        cd ../integration
        
        # Проверяем доступность Docker для интеграционных тестов
        if command -v docker &> /dev/null && docker ps &> /dev/null; then
            go test -v -timeout=20m -tags=integration .
        else
            echo "⚠️  Docker недоступен, пропускаем интеграционные тесты"
            echo "💡 Unit тесты выполнены успешно!"
        fi
        ;;
    *)
        echo "❌ Неверный выбор. Используйте 1, 2 или 3."
        exit 1
        ;;
esac

echo ""
echo "✅ Тесты завершены!"