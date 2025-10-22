package model

import (
	"time"
)

type PartsFilter struct {
	// Список UUID'ов. Пусто — не фильтруем по UUID
	Uuids []string `bson:"uuids,omitempty"`
	// Список имён. Пусто — не фильтруем по имени
	Names []string `bson:"names,omitempty"`
	// Список категорий. Пусто — не фильтруем по категории
	Categories []string `bson:"categories,omitempty"`
	// Список стран производителей. Пусто — не фильтруем по стране
	ManufacturerCountries []string `bson:"manufacturer_countries,omitempty"`
	// Список тегов. Пусто — не фильтруем по тегам
	Tags []string `bson:"tags,omitempty"`
}

type Dimensions struct {
	// Длина в сантиметрах
	Length float64 `bson:"length"`
	// Ширина в сантиметрах
	Width float64 `bson:"width"`
	// Высота в сантиметрах
	Height float64 `bson:"height"`
	// Вес в килограммах
	Weight float64 `bson:"weight"`
}

type Manufacturer struct {
	// Название производителя
	Name string `bson:"name"`
	// Страна производства
	Country string `bson:"country"`
	// Сайт производителя
	Website string `bson:"website"`
}

type Part struct {
	// MongoDB document ID - using string to support UUID format
	ID string `bson:"_id,omitempty"`
	// Уникальный идентификатор детали
	Uuid string `bson:"uuid"`
	// Название детали
	Name string `bson:"name"`
	// Описание детали
	Description string `bson:"description"`
	// Цена за единицу
	Price float64 `bson:"price"`
	// Количество на складе
	StockQuantity int64 `bson:"stock_quantity"`
	// Категория детали
	Category string `bson:"category"`
	// Размеры детали
	Dimensions *Dimensions `bson:"dimensions"`
	// Информация о производителе
	Manufacturer *Manufacturer `bson:"manufacturer"`
	// Теги для быстрого поиска
	Tags []string `bson:"tags"`
	// Гибкие метаданные
	Metadata map[string]interface{} `bson:"metadata"`
	// Дата создания записи
	CreatedAt *time.Time `bson:"created_at"`
	// Дата последнего обновления
	UpdatedAt *time.Time `bson:"updated_at"`
}
