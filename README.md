# 🚀 Rocket Factory
![Coverage](https://img.shields.io/endpoint?url=https://gist.githubusercontent.com/Daniil-Sakharov/10df085c60a8192a207e1efc3d7d3690/raw/coverage.json)

> Микросервисная система для обслуживания и управления ракетами

Современная платформа для заказа, оплаты и управления ракетными компонентами, построенная на микросервисной архитектуре с использованием Go, gRPC и Clean Architecture.

---

## 📋 Содержание

- [О проекте](#о-проекте)
- [Архитектура](#архитектура)
- [Технологии](#технологии)
- [Быстрый старт](#быстрый-старт)
- [Сервисы](#сервисы)
- [API документация](#api-документация)
- [Тестирование](#тестирование)

---

## 🎯 О проекте

**Rocket Factory** — это pet-проект, демонстрирующий best practices разработки микросервисных приложений на Go. Система позволяет:

- 🛒 Создавать заказы на ракетные компоненты
- 💳 Обрабатывать платежи различными способами
- 📦 Управлять складом деталей
- 🔄 Отслеживать статусы заказов

---

## 🏗️ Архитектура

Проект следует принципам **Clean Architecture** и **Domain-Driven Design (DDD)**:

```
┌─────────────────────────────────────────────────────────┐
│                    API Layer (gRPC/HTTP)                │
├─────────────────────────────────────────────────────────┤
│                    Service Layer                        │
│              (Business Logic + Orchestration)           │
├─────────────────────────────────────────────────────────┤
│                  Repository Layer                       │
│                 (Data Access Logic)                     │
├─────────────────────────────────────────────────────────┤
│                   Client Layer                          │
│            (External Service Communication)             │
└─────────────────────────────────────────────────────────┘
```

### Микросервисы

- **Order Service** — управление заказами
- **Payment Service** — обработка платежей
- **Inventory Service** — управление складом деталей

---

## 🛠️ Технологии

### Backend
- **Go 1.24+** — основной язык разработки
- **gRPC** — межсервисное взаимодействие
- **Protocol Buffers** — сериализация данных
- **OpenAPI 3.0** — HTTP API спецификация

### Code Generation
- **Buf** — управление Protobuf схемами
- **ogen** — генерация HTTP серверов из OpenAPI
- **Mockery** — генерация моков для тестирования

### Testing
- **testify** — фреймворк для тестирования
- **testify/suite** — организация тестов
- **gofakeit** — генерация тестовых данных

### Tools
- **Task** — автоматизация задач (аналог Makefile)
- **golangci-lint** — линтинг кода
- **gofumpt** — форматирование кода

---

## 🚀 Быстрый старт

### Требования

- Go 1.24+
- Task (Taskfile)
- Buf CLI

### Установка

```bash
# Клонирование репозитория
git clone https://github.com/yourusername/RocketFactory.git
cd RocketFactory

# Установка зависимостей
go work sync

# Генерация кода (Protobuf, OpenAPI, моки)
task gen

# Запуск всех тестов
task test
```

### Запуск сервисов

```bash
# Order Service (HTTP: 8080, gRPC: 50051)
cd order && go run ./cmd

# Payment Service (gRPC: 50052)
cd payment && go run ./cmd

# Inventory Service (gRPC: 50053)
cd inventory && go run ./cmd
```

---

## 📦 Сервисы

### Order Service

Управление жизненным циклом заказов.

**Возможности:**
- Создание заказа с проверкой наличия деталей
- Оплата заказа через Payment Service
- Получение информации о заказе
- Отмена неоплаченных заказов

**Endpoints:**
- `POST /api/v1/orders` — создать заказ
- `GET /api/v1/orders/{uuid}` — получить заказ
- `POST /api/v1/orders/{uuid}/pay` — оплатить заказ
- `DELETE /api/v1/orders/{uuid}` — отменить заказ

**Swagger UI:** http://localhost:8080/

### Payment Service

Обработка платежей.

**Методы оплаты:**
- 💳 Банковская карта
- 📱 СБП (Система Быстрых Платежей)
- 💰 Кредитная карта
- 🏦 Деньги инвестора

**gRPC API:**
- `PayOrder` — обработать платеж

### Inventory Service

Управление складом ракетных компонентов.

**Категории деталей:**
- 🔧 Двигатели (ENGINE)
- ⛽ Топливо (FUEL)
- 🪟 Иллюминаторы (PORTHOLE)
- ✈️ Крылья (WING)

**gRPC API:**
- `GetPart` — получить деталь по UUID
- `ListParts` — список деталей с фильтрацией 

---

## 📚 API документация

### HTTP API (Order Service)

Swagger UI доступен по адресу: http://localhost:8080/

Автоматическая генерация из OpenAPI спецификации:
```bash
task swagger:generate
```

### gRPC API

Protobuf схемы находятся в `shared/proto/`:
- `inventory/v1/inventory.proto`
- `payment/v1/payment.proto`

Просмотр методов:
```bash
grpcurl -plaintext localhost:50053 list
grpcurl -plaintext localhost:50053 describe inventory.v1.InventoryService
```

---

## 🧪 Тестирование

### Запуск тестов

```bash
# Все тесты
task test

# Конкретный сервис
cd order && go test ./internal/...

# С покрытием
go test ./internal/... -cover
```

### Покрытие тестами

| Сервис | Покрытие |
|--------|----------|
| Order Service | 73.7% |
| Payment Service | 100.0% |
| Inventory Service | 90.9% |

### Структура тестов

Используется паттерн **Test Suite** с `testify/suite`:

```go
type ServiceSuite struct {
    suite.Suite
    ctx             context.Context
    orderRepository *mocks.OrderRepository
    inventoryClient *mocks.InventoryClient
    service         *service
}

func (s *ServiceSuite) TestCreateOrderSuccess() {
    // Arrange
    request := &dto.CreateOrderRequest{...}
    s.inventoryClient.On("ListParts", ...).Return(parts, nil)
    
    // Act
    order, err := s.service.Create(s.ctx, request)
    
    // Assert
    s.Require().NoError(err)
    s.Require().NotNil(order)
}
```

---

## 📁 Структура проекта

```
RocketFactory/
├── order/              # Order Service
│   ├── cmd/           # Entry point
│   ├── internal/      # Внутренняя логика
│   │   ├── api/      # HTTP/gRPC handlers
│   │   ├── service/  # Business logic
│   │   ├── repository/ # Data access
│   │   ├── client/   # External service clients
│   │   └── model/    # Domain models (entity/dto/vo)
│   └── pkg/          # Публичные утилиты
├── payment/           # Payment Service
├── inventory/         # Inventory Service
├── shared/            # Общие компоненты
│   ├── api/          # OpenAPI спецификации
│   ├── proto/        # Protobuf схемы
│   └── pkg/          # Сгенерированный код
├── Taskfile.yml       # Автоматизация задач
└── go.work            # Go workspace
```

---

## 🔧 Доступные команды

```bash
# Генерация кода
task gen                    # Всё (Protobuf + OpenAPI + моки)
task proto:gen              # Только Protobuf
task ogen:gen               # Только OpenAPI
task mocks:gen              # Только моки

# Swagger UI
task swagger:generate       # Генерация Swagger документации
task swagger:clean          # Очистка сгенерированных файлов

# Тестирование
task test                   # Запуск всех тестов

# Линтинг и форматирование
task lint                   # Проверка кода
task fmt                    # Форматирование кода
```

---

## 🎨 Особенности реализации

### Clean Architecture

- **Независимость от фреймворков** — бизнес-логика не зависит от внешних библиотек
- **Тестируемость** — каждый слой тестируется изолированно с использованием моков
- **Разделение ответственности** — четкое разделение на слои (API, Service, Repository, Client)

### Domain-Driven Design

- **Entity** — доменные сущности (Order, Part)
- **Value Objects** — неизменяемые объекты (OrderStatus, PaymentMethod)
- **DTO** — объекты передачи данных между слоями
- **Domain Errors** — централизованное управление ошибками

### Мокирование

Автоматическая генерация моков с помощью Mockery:

```go
s.inventoryClient.On("ListParts", ctx, filter).Return(parts, nil)
s.orderRepository.On("Create", ctx, mock.MatchedBy(func(order *entity.Order) bool {
    return order.UserUUID == userUUID && order.Status == vo.OrderStatusPENDINGPAYMENT
})).Return(nil)
```

---

## 📝 Лицензия

MIT License - свободное использование для обучения и экспериментов.

---

## 🤝 Контакты

Проект создан в образовательных целях для демонстрации микросервисной архитектуры на Go.

**Автор:** Daniil Sakharov  
**GitHub:** [github.com/Daniil-Sakharov](https://github.com/Daniil-Sakharov)

---

⭐ Если проект был полезен, поставьте звезду!