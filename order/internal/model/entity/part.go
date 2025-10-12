package entity

import "time"

// Part - доменная сущность детали (из Inventory)
type Part struct {
	// Уникальный идентификатор детали
	Uuid string
	// Название детали
	Name string
	// Описание детали
	Description string
	// Цена за единицу
	Price float64
	// Количество на складе
	StockQuantity int64
	// Категория детали
	Category Category
	// Размеры детали
	Dimensions *Dimensions
	// Информация о производителе
	Manufacturer *Manufacturer
	// Теги для быстрого поиска
	Tags []string
	// Гибкие метаданные
	Metadata map[string]interface{}
	// Дата создания записи
	CreatedAt *time.Time
	// Дата последнего обновления
	UpdatedAt *time.Time
}

// Category - категория детали
type Category int32

const (
	// Неизвестная категория
	CATEGORY_UNSPECIFIED Category = 0
	// Двигатель
	CATEGORY_ENGINE Category = 1
	// Топливо
	CATEGORY_FUEL Category = 2
	// Иллюминатор
	CATEGORY_PORTHOLE Category = 3
	// Крыло
	CATEGORY_WING Category = 4
)

// Dimensions - размеры детали
type Dimensions struct {
	// Длина в сантиметрах
	Length float64
	// Ширина в сантиметрах
	Width float64
	// Высота в сантиметрах
	Height float64
	// Вес в килограммах
	Weight float64
}

// Manufacturer - информация о производителе
type Manufacturer struct {
	// Название производителя
	Name string
	// Страна производства
	Country string
	// Сайт производителя
	Website string
}

// PartsFilter - фильтр для поиска деталей
type PartsFilter struct {
	// Список UUID'ов. Пусто — не фильтруем по UUID
	Uuids []string
	// Список имён. Пусто — не фильтруем по имени
	Names []string
	// Список категорий. Пусто — не фильтруем по категории
	Categories []Category
	// Список стран производителей. Пусто — не фильтруем по стране
	ManufacturerCountries []string
	// Список тегов. Пусто — не фильтруем по тегам
	Tags []string
}
