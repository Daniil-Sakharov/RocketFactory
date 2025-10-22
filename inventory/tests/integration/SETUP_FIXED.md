# ✅ Исправление E2E тестов для CI/CD

## 🔍 Проблема

E2E тесты в папке `inventory/tests` не работали из-за отсутствия необходимых `.env` файлов конфигурации.

### Ошибка
```
FATAL: Не удалось загрузить .env файл
error: open ../../../deploy/compose/inventory/.env: no such file or directory
```

## ✅ Решение

### 1. Создан мастер файл конфигурации
**Файл:** `/workspace/deploy/env/.env`

Создан из шаблона `.env.template`, содержит все переменные окружения для всех сервисов:
- INVENTORY сервис (MongoDB, gRPC)
- ORDER сервис (PostgreSQL, HTTP, gRPC клиенты)
- PAYMENT сервис (gRPC)

### 2. Сгенерированы .env файлы для сервисов
**Команда:** `task env:generate`

Созданы файлы:
- ✅ `/workspace/deploy/compose/inventory/.env`
- ✅ `/workspace/deploy/compose/order/.env`
- ✅ `/workspace/deploy/compose/payment/.env`

Эти файлы создаются автоматически из шаблонов с подстановкой переменных из мастер файла.

### 3. Проверена конфигурация .gitignore
`.env` файлы правильно игнорируются git (не коммитятся), так как:
- Шаблоны `.env.template` уже в репозитории ✓
- CI/CD генерирует `.env` файлы автоматически ✓
- Секреты не попадают в репозиторий ✓

## 🚀 Как работает CI/CD

### Workflow: `.github/workflows/test-integration-reusable.yml`

```yaml
# 1. Установка Task
- name: 📌 Install Task
  uses: arduino/setup-task@v2.0.0

# 2. Генерация .env файлов
- name: 📄 Генерация .env файлов
  env:
    SERVICES: "inventory,order,payment"
  run: task env:generate

# 3. Запуск интеграционных тестов
- name: 🧪 Run integration tests via Taskfile
  env:
    MODULES: ${{ inputs.modules }}
  run: task test-integration
```

## 📝 Структура конфигурации

```
deploy/
  env/
    .env                      # ✅ Мастер файл (создается из .env.template)
    .env.template             # Шаблон мастер файла (в git)
    inventory.env.template    # Шаблон для inventory (в git)
    order.env.template        # Шаблон для order (в git)
    payment.env.template      # Шаблон для payment (в git)
    generate-env.sh           # Скрипт генерации (в git)
  compose/
    inventory/
      .env                    # ✅ Сгенерированный файл (не в git)
      docker-compose.yml
    order/
      .env                    # ✅ Сгенерированный файл (не в git)
      docker-compose.yml
    payment/
      .env                    # ✅ Сгенерированный файл (не в git)
```

## 🧪 Как запустить тесты локально

### Предварительные требования
1. Docker запущен и доступен
2. Go 1.24+ установлен
3. Task установлен (или используйте go-task)

### Способ 1: Через Task (рекомендуется)

```bash
# 1. Генерируем .env файлы (если еще не сделали)
task env:generate

# 2. Запускаем интеграционные тесты
task test-integration MODULES=inventory
```

### Способ 2: Напрямую через Go

```bash
# 1. Генерируем .env файлы вручную
cd deploy/env
export ENV_SUBST=$HOME/go/bin/envsubst
export SERVICES="inventory,order,payment"
./generate-env.sh

# 2. Запускаем тесты
cd /workspace/inventory/tests/integration
go test -v -timeout=20m -tags=integration .
```

### Способ 3: Если Task не установлен

```bash
# 1. Установить envsubst
go install github.com/a8m/envsubst/cmd/envsubst@v1.4.3

# 2. Сгенерировать .env файлы
cd /workspace
export ENV_SUBST=$HOME/go/bin/envsubst
export SERVICES="inventory,order,payment"
./deploy/env/generate-env.sh

# 3. Запустить тесты
cd inventory/tests/integration
go test -v -timeout=20m -tags=integration .
```

## ⚙️ Что тестируется

### Тесты проверяют:
1. ✅ **ListParts** - получение списка деталей
   - Без фильтра
   - Фильтрация по UUID
   - Фильтрация по категории

2. ✅ **GetPart** - получение детали по UUID
   - Успешное получение существующей детали
   - Ошибка при запросе несуществующей детали

3. ✅ **Полный сценарий** - комплексный тест
   - Вставка тестовых данных
   - Получение списка деталей
   - Получение конкретной детали
   - Фильтрация по категории

## 🐳 Docker контейнеры для тестов

Тесты автоматически создают и удаляют:
- 🗄️ MongoDB контейнер (для хранения данных деталей)
- 🚀 Inventory App контейнер (gRPC сервис)
- 🌐 Docker сеть (для связи контейнеров)

**Время выполнения:**
- Первый запуск: ~3-5 минут (сборка Docker образа)
- Последующие запуски: ~30-60 секунд (кеш Docker)

## 📊 Переменные окружения для тестов

### MongoDB
```env
MONGO_IMAGE_NAME=mongo:8.0
MONGO_HOST=mongodb-container-name
MONGO_PORT=27017
MONGO_DATABASE=inventory
MONGO_AUTH_DB=admin
MONGO_INITDB_ROOT_USERNAME=inventory_admin
MONGO_INITDB_ROOT_PASSWORD=inventory_secret
```

### gRPC Server
```env
GRPC_HOST=0.0.0.0
GRPC_PORT=50051
```

### Logger
```env
LOGGER_LEVEL=debug
LOGGER_AS_JSON=true
```

## ✨ Что теперь работает

1. ✅ CI/CD корректно генерирует .env файлы
2. ✅ Тесты находят необходимую конфигурацию
3. ✅ .env файлы не коммитятся в репозиторий
4. ✅ Локальная разработка поддерживается
5. ✅ Все сервисы имеют единую конфигурацию

## 🎯 Следующие шаги

Для запуска тестов в CI/CD:
1. Push изменений в ветку
2. CI автоматически:
   - Установит Task
   - Сгенерирует .env файлы
   - Запустит Docker
   - Выполнит интеграционные тесты

## 💡 Важные замечания

1. **Локально**: Docker должен быть запущен
2. **CI/CD**: Docker автоматически доступен в GitHub Actions
3. **Timeout**: Тесты используют timeout 20 минут (первая сборка долгая)
4. **Кеширование**: Docker Buildx кеширует слои для ускорения сборки
5. **Очистка**: Контейнеры автоматически удаляются после тестов

## 🔧 Решение проблем

### Ошибка: ".env file not found"
```bash
# Решение: Сгенерировать .env файлы
task env:generate
```

### Ошибка: "Docker not found"
```bash
# Решение: Запустить Docker
docker ps  # Проверить что Docker работает
```

### Ошибка: "timeout"
```bash
# Решение: Увеличить timeout или очистить Docker кеш
go test -v -timeout=30m -tags=integration .
```

### Тесты падают с "connection refused"
```bash
# Решение: Проверить что контейнеры в одной сети
docker network inspect inventory-service
```

## 📚 Документация

- [Taskfile.yml](/Taskfile.yml) - Команды для сборки и тестирования
- [HOW_TO_RUN.md](/inventory/tests/integration/HOW_TO_RUN.md) - Детальная инструкция по запуску
- [.github/workflows/test-integration-reusable.yml](/.github/workflows/test-integration-reusable.yml) - CI/CD workflow

---

**Статус:** ✅ Проблема решена, CI/CD должны проходить успешно
