//go:build integration

package integration

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/joho/godotenv"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.uber.org/zap"
    "os/exec"

	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
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
	err := logger.Init(loggerLevelValue, true)
	if err != nil {
		panic(fmt.Sprintf("не удалось инициализировать логгер: %v", err))
	}

	suiteCtx, suiteCancel = context.WithTimeout(context.Background(), testsTimeout)

    // Готовим окружение Docker для testcontainers и проверяем доступность Docker
    ensureDockerEnv()
    if !isDockerAvailable(suiteCtx) {
        Skip("Docker недоступен в окружении CI — пропускаем интеграционные тесты")
        return
    }

    // Загружаем .env файл и устанавливаем переменные в окружение
    envPath := filepath.Join("..", "..", "..", "deploy", "compose", "inventory", ".env")
    envVars, err := godotenv.Read(envPath)
    if err != nil {
        // В CI/CD файл может отсутствовать — используем безопасные дефолтные значения
        logger.Warn(suiteCtx, "Файл .env не найден, используем значения по умолчанию для тестов",
            zap.String("path", envPath), zap.Error(err))
        envVars = map[string]string{
            // Настройки MongoDB для контейнера
            "MONGO_IMAGE_NAME":             "mongo:8.0",
            "MONGO_INITDB_ROOT_USERNAME":   "inventory-user",
            "MONGO_INITDB_ROOT_PASSWORD":   "inventory-password",
            "MONGO_DATABASE":               "inventory-service",
            "MONGO_AUTH_DB":                "admin",
            // Порт gRPC приложения внутри контейнера (host мапится автоматически)
            "GRPC_PORT":                    "50051",
            // Логгер
            "LOGGER_LEVEL":                 loggerLevelValue,
            "LOGGER_AS_JSON":               "true",
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

// ensureDockerEnv настраивает переменные окружения для testcontainers в CI
func ensureDockerEnv() {
    // Отключаем Ryuk (reaper) и reuse по умолчанию — безопаснее для CI
    setDefaultEnvIfEmpty("TESTCONTAINERS_RYUK_DISABLED", "true")
    setDefaultEnvIfEmpty("TESTCONTAINERS_REUSE_ENABLE", "false")
}

func isDockerAvailable(ctx context.Context) bool {
    // Проверяем, что установлен docker и доступен демон
    if _, err := exec.LookPath("docker"); err != nil {
        logger.Warn(ctx, "docker binary not found in PATH")
        return false
    }
    cmd := exec.Command("docker", "version", "--format", "{{.Server.Version}}")
    if err := cmd.Run(); err != nil {
        logger.Warn(ctx, "docker daemon is not available")
        return false
    }
    return true
}

func setDefaultEnvIfEmpty(key, value string) {
    if os.Getenv(key) == "" {
        _ = os.Setenv(key, value)
    }
}
