package app

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	apiPart "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/api/inventory/v1"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository"
	repoPart "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/part"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service"
	servicePart "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service/part"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

type diContainer struct {
	inventoryV1API      inventoryv1.InventoryServiceServer
	inventoryService    service.PartService
	inventoryRepository repository.PartRepository
	mongoDBClient       *mongo.Client
	mongoDBDatabase     *mongo.Database
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) InventoryAPI(ctx context.Context) inventoryv1.InventoryServiceServer {
	if d.inventoryV1API == nil {
		d.inventoryV1API = apiPart.NewAPI(d.InventoryService(ctx))
	}
	return d.inventoryV1API
}

func (d *diContainer) InventoryService(ctx context.Context) service.PartService {
	if d.inventoryService == nil {
		d.inventoryService = servicePart.NewService(d.InventoryRepository(ctx))
	}
	return d.inventoryService
}

func (d *diContainer) InventoryRepository(ctx context.Context) repository.PartRepository {
	if d.inventoryRepository == nil {
        d.inventoryRepository = repoPart.NewRepository(ctx, d.MongoDBDatabase(ctx))
	}
	return d.inventoryRepository
}

func (d *diContainer) MongoDBClient(ctx context.Context) *mongo.Client {
	if d.mongoDBClient == nil {
		mongoURI := config.AppConfig().Mongo.URI()
		logger.Info(ctx, "Connecting to MongoDB")

		var client *mongo.Client
		var err error

        // Пытаемся подключиться 20 раз с интервалом 3 секунды
        // Используем таймеры, завязанные на ctx, вместо time.Sleep
        maxRetries := 20
        retryDelay := 3 * time.Second

		for attempt := 1; attempt <= maxRetries; attempt++ {
			client, err = mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
			if err != nil {
				logger.Warn(ctx, fmt.Sprintf("Failed to connect to MongoDB (attempt %d/%d): %v", attempt, maxRetries, err))
				time.Sleep(retryDelay)
				continue
			}

			// Проверяем подключение
			err = client.Ping(ctx, readpref.Primary())
			if err == nil {
				logger.Info(ctx, fmt.Sprintf("Successfully connected to MongoDB (attempt %d)", attempt))
				break
			}

			logger.Warn(ctx, fmt.Sprintf("Failed to ping MongoDB (attempt %d/%d): %v", attempt, maxRetries, err))

			// Закрываем неудачное соединение
			_ = client.Disconnect(ctx)
			client = nil

            if attempt < maxRetries {
                timer := time.NewTimer(retryDelay)
                select {
                case <-ctx.Done():
                    timer.Stop()
                case <-timer.C:
                }
            }
		}

		if err != nil || client == nil {
			panic(fmt.Sprintf("failed to connect to MongoDB after %d attempts: %v", maxRetries, err))
		}

		closer.AddNamed("MongoDB client", func(ctx context.Context) error {
			return client.Disconnect(ctx)
		})

		d.mongoDBClient = client
	}

	return d.mongoDBClient
}

func (d *diContainer) MongoDBDatabase(ctx context.Context) *mongo.Database {
	if d.mongoDBDatabase == nil {
		d.mongoDBDatabase = d.MongoDBClient(ctx).Database(config.AppConfig().Mongo.DBName())
	}
	return d.mongoDBDatabase
}
