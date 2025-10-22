# 🔧 CI/CD Setup для Integration Tests

## ✅ Что нужно для работы тестов

### Локальная разработка
1. **Docker** должен быть установлен и запущен
   ```bash
   docker ps  # Проверка, что Docker работает
   ```

2. **Environment файлы** должны быть сгенерированы
   ```bash
   task env:generate
   ```

3. **Запуск тестов**
   ```bash
   task test-integration MODULES=inventory
   # или напрямую:
   go test -v -timeout=20m -tags=integration ./inventory/tests/integration/...
   ```

### CI/CD (GitHub Actions)

Workflow уже настроен в `.github/workflows/test-integration-reusable.yml`:

1. ✅ **Docker Buildx** настроен автоматически
2. ✅ **Environment файлы** генерируются автоматически через `task env:generate`
3. ✅ **Тесты запускаются** через `task test-integration`

## 🚨 Частые проблемы и решения

### Проблема 1: "Не удалось загрузить .env файл"
**Причина:** Не сгенерированы environment файлы

**Решение:**
```bash
task env:generate
```

Или вручную:
```bash
# 1. Создайте базовый .env из шаблона
cp deploy/env/.env.template deploy/env/.env

# 2. Запустите генерацию
cd deploy/env
SERVICES=inventory ./generate-env.sh
```

### Проблема 2: "rootless Docker not found" или "docker: command not found"
**Причина:** Docker не установлен или не запущен

**Решение для локальной разработки:**
```bash
# Проверьте Docker
docker ps

# Если не работает, установите Docker:
# - macOS: Docker Desktop
# - Linux: Docker Engine
# - Windows: Docker Desktop или WSL2 с Docker
```

**Решение для CI/CD:**
- GitHub Actions: Docker уже настроен в workflow (используется `docker/setup-buildx-action@v3`)
- GitLab CI: добавьте `services: - docker:dind` в `.gitlab-ci.yml`
- Jenkins: убедитесь, что Jenkins имеет доступ к Docker socket

### Проблема 3: Тесты падают с timeout
**Причина:** Первый запуск требует сборки Docker образа (~3-5 минут)

**Решение:**
```bash
# Увеличьте timeout
go test -v -timeout=30m -tags=integration ./inventory/tests/integration/...
```

В CI/CD timeout уже установлен в 20 минут в `Taskfile.yml`.

### Проблема 4: "Cannot connect to Docker daemon"
**Причина:** Docker daemon не доступен или нет прав

**Решение:**
```bash
# Linux: добавьте пользователя в группу docker
sudo usermod -aG docker $USER
newgrp docker

# Проверьте права на socket
ls -la /var/run/docker.sock
```

## 📋 Структура тестов

```
inventory/tests/integration/
├── CI_CD_SETUP.md          # Этот файл - инструкции для CI/CD
├── HOW_TO_RUN.md           # Инструкции для локальной разработки
├── suite_test.go           # Настройка test suite (BeforeSuite/AfterSuite)
├── inventory_test.go       # Сами тесты
├── setup.go                # Настройка тестового окружения (Docker контейнеры)
├── teardown.go             # Очистка после тестов
├── test_environment.go     # Вспомогательные методы для работы с БД
└── constants.go            # Константы
```

## 🔍 Как работают тесты

1. **BeforeSuite** (один раз перед всеми тестами):
   - Загружает `.env` файл из `deploy/compose/inventory/.env`
   - Создает Docker сеть
   - Запускает MongoDB в контейнере
   - Собирает и запускает приложение в контейнере
   - Ожидает готовности приложения (wait strategy на gRPC порту)

2. **BeforeEach** (перед каждым тестом):
   - Создает gRPC клиент для подключения к приложению

3. **Тесты**:
   - Вставляют тестовые данные в MongoDB
   - Вызывают gRPC методы
   - Проверяют результаты

4. **AfterEach** (после каждого теста):
   - Очищает коллекцию `parts` в MongoDB

5. **AfterSuite** (один раз после всех тестов):
   - Останавливает контейнер приложения
   - Останавливает контейнер MongoDB
   - Удаляет Docker сеть

## 💡 Советы для CI/CD

### 1. Кеширование Docker образов
В GitHub Actions уже настроено:
```yaml
- name: 🧩 Set up Docker cache
  uses: actions/cache@v4
  with:
    path: /tmp/.buildx-cache
    key: ${{ runner.os }}-buildx-${{ github.sha }}
    restore-keys: |
      ${{ runner.os }}-buildx-
```

### 2. Параллельный запуск тестов
Если у вас несколько модулей, можно запускать тесты параллельно:
```yaml
strategy:
  matrix:
    module: [inventory, order, payment]
```

### 3. Очистка Docker ресурсов
После тестов рекомендуется очистить dangling containers:
```bash
docker container prune -f
```
Уже настроено в GitHub Actions workflow.

## 📊 Метрики производительности

**Первый запуск:**
- Сборка Docker образа: ~3-5 минут
- Запуск MongoDB: ~5 секунд
- Запуск приложения: ~2-6 секунд (retry на подключение к MongoDB)
- Выполнение тестов: ~10-30 секунд
- **Итого: ~4-7 минут**

**Последующие запуски (с кешем):**
- Сборка Docker образа: ~10-30 секунд (кеш слоев)
- Запуск MongoDB: ~5 секунд
- Запуск приложения: ~2-6 секунд
- Выполнение тестов: ~10-30 секунд
- **Итого: ~30-60 секунд**

## 🎯 Чеклист для настройки CI/CD

- [ ] Docker доступен в CI/CD environment
- [ ] Генерация `.env` файлов настроена в pipeline
- [ ] Timeout для тестов установлен >= 20 минут
- [ ] Настроено кеширование Docker образов (опционально)
- [ ] Настроена очистка Docker ресурсов после тестов
- [ ] Проверено, что тесты проходят локально
- [ ] Проверено, что тесты проходят в CI/CD

## 🔗 Полезные ссылки

- [Testcontainers Documentation](https://golang.testcontainers.org/)
- [GitHub Actions - Docker](https://docs.github.com/en/actions/using-containerized-services/about-service-containers)
- [Ginkgo Testing Framework](https://onsi.github.io/ginkgo/)
- [Gomega Matchers](https://onsi.github.io/gomega/)

---

Если возникли проблемы, проверьте логи тестов - они содержат детальную информацию о каждом шаге настройки окружения.
