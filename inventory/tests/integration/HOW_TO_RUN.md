# 🧪 Как запустить E2E тесты для Inventory Service

> 💡 **Для CI/CD setup смотрите [CI_CD_SETUP.md](./CI_CD_SETUP.md)**

## ✅ Что было исправлено

### 1. **Добавлен retry в подключение к MongoDB**
- Файл: `inventory/internal/app/di.go`
- Теперь приложение пытается подключиться к MongoDB **10 раз** с интервалом **2 секунды**
- Это дает MongoDB время для полной инициализации в Docker контейнере

### 2. **Добавлены все необходимые переменные окружения**
- Файл: `inventory/tests/integration/setup.go`
- Контейнер приложения получает все env переменные для подключения к MongoDB

### 3. **Добавлены логи для отладки**
- Теперь видно, какие параметры передаются при запуске контейнера
- Видны попытки подключения к MongoDB

### 4. **Исправлен Dockerfile**
- Версия Go: `1.24.2` → `1.24.5`
- Добавлено копирование `app-inventory` в финальный образ

### 5. **Создан .env файл**
- Путь: `deploy/compose/inventory/.env`
- Содержит настройки для тестов

---

## 🚀 Как запустить тесты

### Способ 1: Через Taskfile (рекомендуется)

```bash
# Убедись, что Docker запущен
docker ps

# Запусти тесты
task test-integration MODULES=inventory
```

### Способ 2: Напрямую через go test

```bash
# Из корня проекта
cd inventory/tests/integration

# Запусти с тегом integration И увеличенным таймаутом
go test -v -timeout=20m -tags=integration .
```

**Важно:** `-timeout=20m` нужен, потому что:
- Сборка Docker образа занимает 3-5 минут
- Стандартный timeout Go test = 10 минут (недостаточно!)
- Наш timeout = 20 минут (с запасом)

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

