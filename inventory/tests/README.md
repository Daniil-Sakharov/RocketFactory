# 🧪 Тесты для Inventory Service

## 📁 Структура тестов

```
inventory/tests/
├── integration/          # Интеграционные тесты (требуют Docker)
│   ├── suite_test.go     # Основной файл с BeforeSuite/AfterSuite
│   ├── inventory_test.go # Тесты gRPC API
│   ├── setup.go          # Настройка тестового окружения
│   ├── test_environment.go # Вспомогательные функции
│   ├── teardown.go       # Очистка после тестов
│   ├── constants.go      # Константы
│   ├── simple_test.go    # Простые тесты
│   ├── mock_test.go      # Мок-тесты
│   └── standalone_test.go # Standalone тесты
└── unit/                 # Unit тесты (без внешних зависимостей)
    └── inventory_unit_test.go
```

## 🚀 Как запустить тесты

### Unit тесты (рекомендуется для CI/CD)

```bash
# Из корня проекта
cd inventory/tests/unit
go test -v -timeout=20m -tags=unit .
```

### Интеграционные тесты (требуют Docker)

```bash
# Из корня проекта
cd inventory/tests/integration
go test -v -timeout=20m -tags=integration .
```

### Все тесты сразу

```bash
# Unit тесты
go test -v -timeout=20m -tags=unit ./inventory/tests/unit/...

# Интеграционные тесты (если Docker доступен)
go test -v -timeout=20m -tags=integration ./inventory/tests/integration/...
```

## 🔧 Требования

### Unit тесты
- ✅ Go 1.24+
- ✅ Зависимости из go.mod

### Интеграционные тесты
- ✅ Go 1.24+
- ✅ Docker (запущен и доступен)
- ✅ Переменные окружения в `.env` файле

## 📋 Что тестируется

### Unit тесты
- ✅ Валидация структур данных
- ✅ Бизнес-логика
- ✅ Конвертация данных
- ✅ Фильтрация
- ✅ Обработка ошибок

### Интеграционные тесты
- ✅ gRPC API endpoints
- ✅ Подключение к MongoDB
- ✅ Полный цикл работы с данными
- ✅ Docker контейнеры

## 🐛 Устранение проблем

### Docker недоступен
Если Docker не запущен или недоступен, используйте unit тесты:
```bash
go test -v -timeout=20m -tags=unit ./inventory/tests/unit/...
```

### Отсутствует .env файл
Создайте .env файл:
```bash
cp deploy/env/inventory.env.template deploy/compose/inventory/.env
```

### Проблемы с правами Docker
```bash
sudo usermod -aG docker $USER
# Перезайдите в систему
```

## 📊 Результаты тестов

### Успешный запуск unit тестов
```
=== RUN   TestInventoryUnit
Running Suite: Inventory Service Unit Test Suite
==================================================================================
Random Seed: 1761139488

Will run 7 of 7 specs
•••••••

Ran 7 of 7 Specs in 0.000 seconds
SUCCESS! -- 7 Passed | 0 Failed | 0 Pending | 0 Skipped
--- PASS: TestInventoryUnit (0.00s)
PASS
ok  	github.com/Daniil-Sakharov/RocketFactory/inventory/tests/unit	0.005s
```

## 🎯 Для CI/CD

Рекомендуется использовать unit тесты в CI/CD пайплайне, так как они:
- ✅ Не требуют Docker
- ✅ Быстро выполняются
- ✅ Не зависят от внешних сервисов
- ✅ Покрывают основную бизнес-логику

```yaml
# Пример для GitHub Actions
- name: Run Unit Tests
  run: |
    cd inventory/tests/unit
    go test -v -timeout=20m -tags=unit .
```