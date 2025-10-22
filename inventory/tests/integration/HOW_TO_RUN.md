# 🧪 Как запустить E2E тесты для Inventory Service

## ✅ Что было исправлено для работы в CI/CD

### 1. **Создан главный файл конфигурации**
- Файл: `deploy/env/.env`
- Содержит все переменные окружения для всех сервисов проекта
- Используется скриптом генерации для создания service-specific .env файлов

### 2. **Создан .env файл для Inventory тестов**
- Путь: `deploy/compose/inventory/.env`
- Генерируется автоматически из шаблона `deploy/env/inventory.env.template`
- Содержит все необходимые настройки для MongoDB и gRPC

### 3. **MongoDB retry уже настроен**
- Файл: `inventory/internal/app/di.go`
- Приложение пытается подключиться к MongoDB **20 раз** с интервалом **3 секунды**
- Это дает MongoDB достаточно времени для инициализации в Docker контейнере (до 60 секунд)

### 4. **Настроен GitHub Workflow**
- Workflow автоматически генерирует .env файлы через `task env:generate`
- Запускает тесты с Docker через testcontainers
- Очищает контейнеры после завершения тестов

### 5. **Конфигурация переменных окружения**
Все необходимые переменные для тестов:
- `GRPC_HOST`, `GRPC_PORT` - для gRPC сервера
- `LOGGER_LEVEL`, `LOGGER_AS_JSON` - для логирования
- `MONGO_IMAGE_NAME` - образ MongoDB для testcontainers
- `MONGO_HOST`, `MONGO_PORT`, `MONGO_DATABASE` - подключение к MongoDB
- `MONGO_AUTH_DB`, `MONGO_INITDB_ROOT_USERNAME`, `MONGO_INITDB_ROOT_PASSWORD` - аутентификация

---

## 🚀 Как запустить тесты

### Способ 1: Через GitHub Actions (автоматически в CI/CD)

Тесты запускаются автоматически при push или создании pull request:

```yaml
# Workflow выполняет:
1. Генерацию .env файлов: task env:generate
2. Запуск тестов: task test-integration MODULES=inventory
3. Очистку Docker контейнеров
```

### Способ 2: Локально через Taskfile (рекомендуется)

```bash
# 1. Убедись, что Docker запущен
docker ps

# 2. Сгенерируй .env файлы (один раз)
task env:generate

# 3. Запусти тесты
task test-integration MODULES=inventory
```

### Способ 3: Напрямую через go test

```bash
# 1. Убедись что .env файл существует
ls deploy/compose/inventory/.env

# 2. Из корня проекта, запусти с тегом integration
cd inventory/tests/integration
go test -v -timeout=20m -tags=integration .
```

**Важно:** 
- `-timeout=20m` нужен для сборки Docker образа (3-5 минут при первом запуске)
- Тесты требуют Docker для создания testcontainers
- .env файл должен существовать в `deploy/compose/inventory/.env`

---

## 🔍 Если тесты все еще падают

### 1. Проверь, что Docker запущен

```bash
docker ps
# Должен вернуть список контейнеров (может быть пустой)
```

### 2. Проверь логи во время выполнения теста

В **другом терминале** во время выполнения тестов:

```bash
# Смотри контейнеры
docker ps --filter "label=org.testcontainers=true"

# Смотри логи MongoDB
docker logs <mongo_container_id>

# Смотри логи приложения
docker logs <app_container_id>
```

### 3. Проверь, что контейнеры в одной сети

```bash
# Найди сеть
docker network ls | grep inventory

# Проверь контейнеры в сети
docker network inspect <network_name>
```

### 4. Вручную проверь, что приложение может подключиться к MongoDB

```bash
# 1. Создай сеть
docker network create test-network

# 2. Запусти MongoDB
docker run -d --name test-mongo \
  --network test-network \
  -e MONGO_INITDB_ROOT_USERNAME=inventory-user \
  -e MONGO_INITDB_ROOT_PASSWORD=inventory-password \
  mongo:8.0

# 3. Запусти приложение
docker run --rm \
  --network test-network \
  -e GRPC_PORT=50051 \
  -e MONGO_HOST=test-mongo \
  -e MONGO_PORT=27017 \
  -e MONGO_DATABASE=inventory-service \
  -e MONGO_INITDB_ROOT_USERNAME=inventory-user \
  -e MONGO_INITDB_ROOT_PASSWORD=inventory-password \
  -e MONGO_AUTH_DB=admin \
  -e LOGGER_LEVEL=debug \
  -e LOGGER_AS_JSON=true \
  test-inventory:latest

# Приложение должно стартовать и подключиться к MongoDB через 1-2 попытки
# В логах увидишь: "Successfully connected to MongoDB (attempt X)"
```

### 5. Очисти Docker кеш (если ничего не помогает)

```bash
# Останови все testcontainers
docker ps -a --filter "label=org.testcontainers=true" --format "{{.ID}}" | xargs -r docker rm -f

# Удали старые образы
docker images | grep inventory | awk '{print $3}' | xargs -r docker rmi -f

# Пересобери образ
docker build -f deploy/docker/inventory/Dockerfile -t test-inventory:latest .
```

---

## 📊 Ожидаемый результат

После исправлений тесты должны:

1. ✅ **Создать Docker сеть** (~1 секунда)
2. ✅ **Запустить MongoDB контейнер** (~5 секунд)
3. ✅ **Собрать Docker образ приложения** (~3-5 минут при первом запуске, потом быстрее благодаря кешу)
4. ✅ **Запустить контейнер приложения** (~2 секунды)
5. ✅ **Приложение подключается к MongoDB** (1-3 попытки, ~2-6 секунд)
6. ✅ **Запустить тесты** (~10-30 секунд)
7. ✅ **Очистить все контейнеры** (~2 секунды)

**Общее время**: 4-7 минут при первом запуске, 30-60 секунд при последующих (благодаря кешу Docker).

---

## 💡 Подсказки

- **Первый запуск всегда долгий** — Docker скачивает образы и собирает приложение
- **Последующие запуски быстрые** — используется кеш
- **Если таймаут**: увеличь `startupTimeout` в `setup.go` (сейчас 5 минут)
- **Если MongoDB не успевает стартовать**: увеличь `maxRetries` в `di.go` (сейчас 10 попыток)

---

## 🎯 Контрольные вопросы для проверки понимания

1. **Почему нужен retry при подключении к MongoDB?**
2. **Что произойдет, если контейнеры окажутся в разных Docker сетях?**
3. **Зачем нужен build tag `//go:build integration`?**

Удачи! 🚀

