package part

import (
	"time"

	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
)

func (r *repository) InitTestData() {
	r.mu.Lock()
	defer r.mu.Unlock()

	now := time.Now()

	// Тестовые данные
	testParts := []repoModel.Part{
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Ракетный двигатель RD-180",
			Description:   "Мощный жидкостный ракетный двигатель",
			Price:         15000000.0,
			StockQuantity: 3,
			Category:      repoModel.CATEGORY_ENGINE,
			Dimensions: &repoModel.Dimensions{
				Length: 350.0,
				Width:  240.0,
				Height: 240.0,
				Weight: 5480.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Энергомаш",
				Country: "Russia",
				Website: "www.npoenergomash.ru",
			},
			Tags: []string{"двигатель", "мощный", "жидкостный"},
			Metadata: map[string]interface{}{
				"тяга":    3827000.0,
				"топливо": "керосин+кислород",
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Крыло Delta-V",
			Description:   "Аэродинамическое крыло для атмосферного полета",
			Price:         2500000.0,
			StockQuantity: 8,
			Category:      repoModel.CATEGORY_WING,
			Dimensions: &repoModel.Dimensions{
				Length: 1200.0,
				Width:  600.0,
				Height: 50.0,
				Weight: 850.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"крыло", "аэродинамика", "композит"},
			Metadata: map[string]interface{}{
				"материал":       "углеродное волокно",
				"термостойкость": 1500.0,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Криогенное топливо LH2",
			Description:   "Жидкий водород для ракетных двигателей",
			Price:         50000.0,
			StockQuantity: 150,
			Category:      repoModel.CATEGORY_FUEL,
			Dimensions: &repoModel.Dimensions{
				Length: 100.0,
				Width:  100.0,
				Height: 200.0,
				Weight: 70.8,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Air Liquide",
				Country: "France",
				Website: "www.airliquide.com",
			},
			Tags: []string{"топливо", "криогенное", "водород"},
			Metadata: map[string]interface{}{
				"температура": -253.0,
				"чистота":     99.9,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Иллюминатор Космос-360",
			Description:   "Прочный иллюминатор для наблюдения в космосе",
			Price:         750000.0,
			StockQuantity: 12,
			Category:      repoModel.CATEGORY_PORTHOLE,
			Dimensions: &repoModel.Dimensions{
				Length: 60.0,
				Width:  60.0,
				Height: 15.0,
				Weight: 25.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Roscosmos",
				Country: "Russia",
				Website: "www.roscosmos.ru",
			},
			Tags: []string{"иллюминатор", "обзор", "прочный"},
			Metadata: map[string]interface{}{
				"материал_стекла": "сапфировое стекло",
				"давление":        101325.0,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Двигатель Merlin 1D",
			Description:   "Компактный двигатель для первой ступени",
			Price:         1200000.0,
			StockQuantity: 25,
			Category:      repoModel.CATEGORY_ENGINE,
			Dimensions: &repoModel.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 630.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"двигатель", "компактный", "многоразовый"},
			Metadata: map[string]interface{}{
				"тяга":         845000.0,
				"многоразовый": true,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440006",
			Name:          "Крыло Falcon Heavy",
			Description:   "Большое крыло для тяжелых ракет",
			Price:         4200000.0,
			StockQuantity: 4,
			Category:      repoModel.CATEGORY_WING,
			Dimensions: &repoModel.Dimensions{
				Length: 1800.0,
				Width:  900.0,
				Height: 80.0,
				Weight: 1500.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Blue Origin",
				Country: "USA",
				Website: "www.blueorigin.com",
			},
			Tags: []string{"крыло", "тяжелое", "стабилизация"},
			Metadata: map[string]interface{}{
				"грузоподъемность": 63800.0,
				"сертификат":       "NASA-2024",
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440007",
			Name:          "Топливный бак LOX-5000",
			Description:   "Бак для жидкого кислорода большой емкости",
			Price:         890000.0,
			StockQuantity: 18,
			Category:      repoModel.CATEGORY_FUEL,
			Dimensions: &repoModel.Dimensions{
				Length: 500.0,
				Width:  200.0,
				Height: 200.0,
				Weight: 1200.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "Airbus Defence and Space",
				Country: "Germany",
				Website: "www.airbus.com",
			},
			Tags: []string{"топливо", "бак", "кислород"},
			Metadata: map[string]interface{}{
				"объем":    5000.0,
				"давление": 300.0,
				"изоляция": "криогенная",
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440008",
			Name:          "Иллюминатор Starship View",
			Description:   "Панорамный иллюминатор для туристических полетов",
			Price:         1250000.0,
			StockQuantity: 6,
			Category:      repoModel.CATEGORY_PORTHOLE,
			Dimensions: &repoModel.Dimensions{
				Length: 150.0,
				Width:  100.0,
				Height: 20.0,
				Weight: 85.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"иллюминатор", "панорамный", "туризм"},
			Metadata: map[string]interface{}{
				"угол_обзора": 180.0,
				"UV_защита":   true,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440009",
			Name:          "Raptor Engine V2",
			Description:   "Полнопоточный двигатель на метане",
			Price:         2800000.0,
			StockQuantity: 15,
			Category:      repoModel.CATEGORY_ENGINE,
			Dimensions: &repoModel.Dimensions{
				Length: 340.0,
				Width:  130.0,
				Height: 130.0,
				Weight: 1600.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"двигатель", "метан", "полнопоточный"},
			Metadata: map[string]interface{}{
				"тяга":             2300000.0,
				"удельный_импульс": 330.0,
				"топливо":          "метан+кислород",
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
		{
			Uuid:          "550e8400-e29b-41d4-a716-446655440010",
			Name:          "Grid Fins",
			Description:   "Решетчатые рули для управления посадкой",
			Price:         650000.0,
			StockQuantity: 20,
			Category:      repoModel.CATEGORY_WING,
			Dimensions: &repoModel.Dimensions{
				Length: 150.0,
				Width:  120.0,
				Height: 50.0,
				Weight: 180.0,
			},
			Manufacturer: &repoModel.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"крыло", "управление", "посадка"},
			Metadata: map[string]interface{}{
				"материал":     "титан",
				"многоразовый": true,
			},
			CreatedAt: &now,
			UpdatedAt: &now,
		},
	}

	// Загружаем данные в мапу через цикл
	for _, part := range testParts {
		r.data[part.Uuid] = part
	}
}
