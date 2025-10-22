//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
    "os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
    "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/testcontainers"
)

const testsTimeout = 60 * time.Minute // Увеличен для сборки Docker образа (каждый раз собирается с нуля)

var (
	env *TestEnvironment

	suiteCtx    context.Context
	suiteCancel context.CancelFunc
)

func TestIntegration(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Inventory Service Integration Test Suite")
}

var _ = BeforeSuite(func() {
    // Пропускаем весь набор интеграционных тестов, если Docker недоступен
    if _, err := exec.LookPath("docker"); err != nil {
        Skip("Docker не найден в системе — пропускаем интеграционные тесты")
    }
    if err := exec.Command("docker", "ps").Run(); err != nil {
        Skip("Docker демон недоступен — пропускаем интеграционные тесты")
    }

	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

	// Загружаем .env файл и устанавливаем переменные в окружение
    envVars, err := godotenv.Read(filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env"))
    if err != nil {
        logger.Warn(suiteCtx, "Не удалось загрузить .env файл — используем значения по умолчанию", zap.Error(err))
        envVars = map[string]string{
            // Mongo defaults
            testcontainers.MongoImageNameKey: "mongo:8.0",
            testcontainers.MongoUsernameKey:   "inventory-user",
            testcontainers.MongoPasswordKey:   "inventory-password",
            testcontainers.MongoDatabaseKey:   "inventory-service",
            testcontainers.MongoAuthDBKey:     "admin",
            // App defaults
            grpcPortKey:    "50051",
            "LOGGER_LEVEL":   loggerLevelValue,
            "LOGGER_AS_JSON": "true",
        }
    }

	// Устанавливаем переменные в окружение процесса
	for key, value := range envVars {
		_ = os.Setenv(key, value)
	}

	logger.Info(suiteCtx, "Запуск тестового окружения...")
	env = setupTestEnvironment(suiteCtx)
})

var _ = AfterSuite(func() {
	logger.Info(context.Background(), "Завершение набора тестов")
	if env != nil {
		teardownTestEnvironment(suiteCtx, env)
	}
	suiteCancel()
})
