//go:build integration

package integration

import (
	"context"
	"os"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers/app"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers/mongo"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers/network"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers/path"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.uber.org/zap"
)

const (
	// Параметры для контейнеров
	inventoryAppName    = "inventory-app"
	inventoryDockerfile = "deploy/docker/inventory/Dockerfile"

	// Переменные окружения приложения
	grpcHostKey = "GRPC_HOST"
	grpcPortKey = "GRPC_PORT"

	// Значения переменных окружения
	grpcHostValue    = "0.0.0.0"
	loggerLevelValue = "debug"
	startupTimeout   = 45 * time.Minute // Увеличен для сборки Docker образа (go mod download занимает ~15 минут)
)

type TestEnvironment struct {
	Network *network.Network
	Mongo   *mongo.Container
	App     *app.Container
}

func setupTestEnvironment(ctx context.Context) *TestEnvironment {
	logger.Info(ctx, "🚀 Подготовка тестового окружения...")

	generatedNetwork, err := network.NewNetwork(ctx, "inventory-service")
	if err != nil {
		logger.Fatal(ctx, "не удалось создать общую сеть", zap.Error(err))
	}
	logger.Info(ctx, "✅ Сеть успешно создана")

	mongoUsername := getEnvWithLogging(ctx, "MONGO_INITDB_ROOT_USERNAME")
	mongoPassword := getEnvWithLogging(ctx, "MONGO_INITDB_ROOT_PASSWORD")
	mongoImageName := getEnvWithLogging(ctx, "MONGO_IMAGE_NAME")
	mongoDatabase := getEnvWithLogging(ctx, "MONGO_DATABASE")

	grpcPort := getEnvWithLogging(ctx, "GRPC_PORT")

	generatedMongo, err := mongo.NewContainer(ctx,
		mongo.WithNetworkName(generatedNetwork.Name()),
		mongo.WithContainerName(testcontainers.MongoContainerName),
		mongo.WithImageName(mongoImageName),
		mongo.WithDatabase(mongoDatabase),
		mongo.WithAuth(mongoUsername, mongoPassword),
		mongo.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork})
		logger.Fatal(ctx, "не удалось запустить контейнер MongoDB", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер MongoDB успешно запущен")

	projectRoot := path.GetProjectRoot()

	appEnv := map[string]string{
		// Переопределяем хост MongoDB для подключения к контейнеру из testcontainers
		"MONGO_HOST":     generatedMongo.Config().ContainerName,
		"MONGO_PORT":     "27017",
		"MONGO_DATABASE": mongoDatabase,
		"MONGO_INITDB_ROOT_USERNAME": mongoUsername,
		"MONGO_INITDB_ROOT_PASSWORD": mongoPassword,
		"MONGO_AUTH_DB":   getEnvWithLogging(ctx, "MONGO_AUTH_DB"),
		"GRPC_HOST":                     grpcHostValue,
		"GRPC_PORT":                     grpcPort,
		"LOGGER_LEVEL":                  "debug",
		"LOGGER_AS_JSON":                "true",
	}


	logger.Info(ctx, "🚀 Starting app container",
		zap.String("network", generatedNetwork.Name()),
		zap.String("mongo_host", generatedMongo.Config().ContainerName),
		zap.String("grpc_host", grpcHostValue),
		zap.String("grpc_port", grpcPort),
		zap.Any("env", appEnv))

	waitStrategy := wait.ForListeningPort(nat.Port(grpcPort + "/tcp")).
		WithStartupTimeout(startupTimeout)

	appContainer, err := app.NewContainer(ctx,
		app.WithName(inventoryAppName),
		app.WithPort(grpcPort),
		app.WithDockerfile(projectRoot, inventoryDockerfile),
		app.WithNetwork(generatedNetwork.Name()),
		app.WithEnv(appEnv),
		app.WithLogOutput(os.Stdout),
		app.WithStartupWait(waitStrategy),
		app.WithLogger(logger.Logger()),
	)
	if err != nil {
		cleanupTestEnvironment(ctx, &TestEnvironment{Network: generatedNetwork, Mongo: generatedMongo})
		logger.Fatal(ctx, "не удалось запустить контейнер приложения", zap.Error(err))
	}
	logger.Info(ctx, "✅ Контейнер приложения успешно запущен")

	logger.Info(ctx, "🎉 Тестовое окружение готово")
	return &TestEnvironment{
		Network: generatedNetwork,
		Mongo:   generatedMongo,
		App:     appContainer,
	}
}

func getEnvWithLogging(ctx context.Context, key string) string {
	value := os.Getenv(key)
	if value == "" {
		logger.Warn(ctx, "Переменная окружения не установлена", zap.String("key", key))
	}

	return value
}
