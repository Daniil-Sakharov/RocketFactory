//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/converter"
	repoConverter "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

// InsertTestPart — вставляет одну тестовую деталь в коллекцию Mongo и возвращает её UUID
func (env *TestEnvironment) InsertTestPart(ctx context.Context) (string, error) {
	partUUID := gofakeit.UUID()
	now := time.Now()

	partDoc := bson.M{
		"_id":            partUUID,
		"uuid":           partUUID, // ВАЖНО: добавляем поле uuid для корректного поиска
		"name":           gofakeit.CarMaker() + " " + gofakeit.CarModel(),
		"description":    gofakeit.Sentence(10),
		"price":          gofakeit.Float64Range(100.0, 100000.0),
		"stock_quantity": int64(gofakeit.Number(0, 1000)),
		"category":       "ENGINE",
		"dimensions": bson.M{
			"length": gofakeit.Float64Range(10.0, 500.0),
			"width":  gofakeit.Float64Range(10.0, 500.0),
			"height": gofakeit.Float64Range(10.0, 500.0),
			"weight": gofakeit.Float64Range(5.0, 5000.0),
		},
		"manufacturer": bson.M{
			"name":    gofakeit.Company(),
			"country": gofakeit.Country(),
			"website": gofakeit.URL(),
		},
		"tags": []string{
			gofakeit.Word(),
			gofakeit.Word(),
			gofakeit.Word(),
		},
		"metadata": bson.M{
			"quality":        "premium",
			"warranty_years": int64(5),
		},
		"created_at": primitive.NewDateTimeFromTime(now),
		"updated_at": primitive.NewDateTimeFromTime(now),
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение должно совпадать с .env
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// InsertTestParts — вставляет несколько тестовых деталей в коллекцию Mongo
func (env *TestEnvironment) InsertTestParts(ctx context.Context) error {
	now := time.Now()

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение должно совпадать с .env
	}

	// Log database name for debugging
	println("DEBUG: Inserting test data into database:", databaseName, "collection:", partsCollectionName) //nolint:forbidigo // Debug logging for tests

	// Создаем 5 разных деталей с разными категориями
	uuid1 := gofakeit.UUID()
	uuid2 := gofakeit.UUID()
	uuid3 := gofakeit.UUID()
	uuid4 := gofakeit.UUID()
	uuid5 := gofakeit.UUID()

	parts := []interface{}{
		// Двигатель
		bson.M{
			"_id":            uuid1,
			"uuid":           uuid1, // ВАЖНО: добавляем поле uuid для корректного поиска
			"name":           "Ионный двигатель X-500",
			"description":    "Высокоэффективный ионный двигатель для межпланетных перелетов",
			"price":          50000.0,
			"stock_quantity": int64(10),
			"category":       "ENGINE",
			"dimensions": bson.M{
				"length": 250.0,
				"width":  150.0,
				"height": 200.0,
				"weight": 500.0,
			},
			"manufacturer": bson.M{
				"name":    "SpaceTech Industries",
				"country": "USA",
				"website": "https://spacetech.com",
			},
			"tags":       []string{"ion", "efficient", "space"},
			"metadata":   bson.M{"power_output": int64(5000), "fuel_type": "xenon"},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		// Топливо
		bson.M{
			"_id":            uuid2,
			"uuid":           uuid2,
			"name":           "Жидкий ксенон премиум",
			"description":    "Высокочистый жидкий ксенон для ионных двигателей",
			"price":          1500.0,
			"stock_quantity": int64(500),
			"category":       "FUEL",
			"dimensions": bson.M{
				"length": 50.0,
				"width":  30.0,
				"height": 30.0,
				"weight": 25.0,
			},
			"manufacturer": bson.M{
				"name":    "Quantum Fuels",
				"country": "Germany",
				"website": "https://quantumfuels.de",
			},
			"tags":       []string{"xenon", "fuel", "premium"},
			"metadata":   bson.M{"purity": "99.999%", "volume_liters": int64(10)},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		// Иллюминатор
		bson.M{
			"_id":            uuid3,
			"uuid":           uuid3,
			"name":           "Панорамный иллюминатор AstroView",
			"description":    "Многослойный защищенный иллюминатор с антибликовым покрытием",
			"price":          8000.0,
			"stock_quantity": int64(25),
			"category":       "PORTHOLE",
			"dimensions": bson.M{
				"length": 80.0,
				"width":  80.0,
				"height": 15.0,
				"weight": 45.0,
			},
			"manufacturer": bson.M{
				"name":    "ClearView Corp",
				"country": "Japan",
				"website": "https://clearview.jp",
			},
			"tags":       []string{"porthole", "panoramic", "protected"},
			"metadata":   bson.M{"layers": int64(5), "uv_protection": true},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		// Крыло
		bson.M{
			"_id":            uuid4,
			"uuid":           uuid4,
			"name":           "Аэродинамическое крыло Delta-9",
			"description":    "Титановое крыло с регулируемой геометрией для атмосферного полета",
			"price":          35000.0,
			"stock_quantity": int64(5),
			"category":       "WING",
			"dimensions": bson.M{
				"length": 800.0,
				"width":  300.0,
				"height": 50.0,
				"weight": 1200.0,
			},
			"manufacturer": bson.M{
				"name":    "AeroDynamics Ltd",
				"country": "France",
				"website": "https://aerodynamics.fr",
			},
			"tags":       []string{"wing", "titanium", "adjustable"},
			"metadata":   bson.M{"material": "titanium-alloy", "max_speed_kmh": int64(25000)},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
		// Еще один двигатель
		bson.M{
			"_id":            uuid5,
			"uuid":           uuid5,
			"name":           "Плазменный двигатель Nebula-7",
			"description":    "Компактный плазменный двигатель нового поколения",
			"price":          75000.0,
			"stock_quantity": int64(3),
			"category":       "ENGINE",
			"dimensions": bson.M{
				"length": 180.0,
				"width":  120.0,
				"height": 150.0,
				"weight": 350.0,
			},
			"manufacturer": bson.M{
				"name":    "Nebula Propulsion",
				"country": "Russia",
				"website": "https://nebulaprop.ru",
			},
			"tags":       []string{"plasma", "compact", "next-gen"},
			"metadata":   bson.M{"power_output": int64(7500), "fuel_type": "plasma"},
			"created_at": primitive.NewDateTimeFromTime(now),
			"updated_at": primitive.NewDateTimeFromTime(now),
		},
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertMany(ctx, parts)
	if err != nil {
		return err
	}

	return nil
}

// InsertTestPartWithData — вставляет тестовую деталь с заданными данными
func (env *TestEnvironment) InsertTestPartWithData(ctx context.Context, part *inventoryV1.Part) (string, error) {
	partUUID := part.Uuid
	if partUUID == "" {
		partUUID = gofakeit.UUID()
	}

	now := time.Now()
	createdAt := now
	updatedAt := now

	if part.CreatedAt != nil {
		createdAt = part.CreatedAt.AsTime()
	}
	if part.UpdatedAt != nil {
		updatedAt = part.UpdatedAt.AsTime()
	}

	// Конвертируем protobuf Category в строку для MongoDB
	categoryStr := repoConverter.CategoryToRepoModel(converter.CategoryFromProto(part.Category))

	partDoc := bson.M{
		"_id":            partUUID,
		"uuid":           partUUID, // ВАЖНО: добавляем поле uuid
		"name":           part.Name,
		"description":    part.Description,
		"price":          part.Price,
		"stock_quantity": part.StockQuantity,
		"category":       categoryStr,
		"created_at":     primitive.NewDateTimeFromTime(createdAt),
		"updated_at":     primitive.NewDateTimeFromTime(updatedAt),
	}

	// Добавляем dimensions если есть
	if part.Dimensions != nil {
		partDoc["dimensions"] = bson.M{
			"length": part.Dimensions.Length,
			"width":  part.Dimensions.Width,
			"height": part.Dimensions.Height,
			"weight": part.Dimensions.Weight,
		}
	}

	// Добавляем manufacturer если есть
	if part.Manufacturer != nil {
		partDoc["manufacturer"] = bson.M{
			"name":    part.Manufacturer.Name,
			"country": part.Manufacturer.Country,
			"website": part.Manufacturer.Website,
		}
	}

	// Добавляем tags если есть
	if len(part.Tags) > 0 {
		partDoc["tags"] = part.Tags
	}

	// Добавляем metadata если есть
	if len(part.Metadata) > 0 {
		metadata := bson.M{}
		for key, value := range part.Metadata {
			switch v := value.Value.(type) {
			case *inventoryV1.Value_StringValue:
				metadata[key] = v.StringValue
			case *inventoryV1.Value_Int64Value:
				metadata[key] = v.Int64Value
			case *inventoryV1.Value_DoubleValue:
				metadata[key] = v.DoubleValue
			case *inventoryV1.Value_BoolValue:
				metadata[key] = v.BoolValue
			}
		}
		partDoc["metadata"] = metadata
	}

	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение должно совпадать с .env
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).InsertOne(ctx, partDoc)
	if err != nil {
		return "", err
	}

	return partUUID, nil
}

// GetTestPartData — возвращает тестовые данные детали
func (env *TestEnvironment) GetTestPartData() *inventoryV1.Part {
	return &inventoryV1.Part{
		Uuid:          gofakeit.UUID(),
		Name:          "Тестовый двигатель TestEngine-1",
		Description:   "Это тестовая деталь для проверки функциональности",
		Price:         99999.99,
		StockQuantity: 100,
		Category:      inventoryV1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryV1.Dimensions{
			Length: 200.0,
			Width:  100.0,
			Height: 150.0,
			Weight: 450.0,
		},
		Manufacturer: &inventoryV1.Manufacturer{
			Name:    "Test Manufacturing Co.",
			Country: "TestLand",
			Website: "https://testmfg.test",
		},
		Tags: []string{"test", "engine", "sample"},
		Metadata: map[string]*inventoryV1.Value{
			"test_mode": {
				Value: &inventoryV1.Value_BoolValue{BoolValue: true},
			},
			"test_iteration": {
				Value: &inventoryV1.Value_Int64Value{Int64Value: 1},
			},
		},
		CreatedAt: timestamppb.New(time.Now()),
		UpdatedAt: timestamppb.New(time.Now()),
	}
}

// ClearPartsCollection — удаляет все записи из коллекции parts
func (env *TestEnvironment) ClearPartsCollection(ctx context.Context) error {
	// Используем базу данных из переменной окружения MONGO_DATABASE
	databaseName := os.Getenv("MONGO_DATABASE")
	if databaseName == "" {
		databaseName = "inventory" // fallback значение должно совпадать с .env
	}

	_, err := env.Mongo.Client().Database(databaseName).Collection(partsCollectionName).DeleteMany(ctx, bson.M{})
	if err != nil {
		return err
	}

	return nil
}
