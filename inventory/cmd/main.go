package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/app"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	"go.uber.org/zap"
)

func main() {
	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	// Перехватываем панику для отладки
	defer func() {
		if r := recover(); r != nil {
			logger.Error(appCtx, "🔥 PANIC occurred", zap.Any("panic", r))
			panic(r) // Повторно бросаем панику
		}
	}()

	// .env файл опционален:
	// - В локальной разработке: загружается из корня проекта или указанного пути
	// - В Docker: переменные передаются через environment (-e флаги)
	err := config.Load()
	if err != nil {
		logger.Error(appCtx, "❌ Failed to load config", zap.Error(err))
		panic(fmt.Errorf("error to load config: %w", err))
	}
	logger.Info(appCtx, "✅ Config loaded")

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	logger.Info(appCtx, "🏗️ Creating application...")
	a, err := app.New(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ Failed to create app", zap.Error(err))
		return
	}
	logger.Info(appCtx, "✅ Application created")

	logger.Info(appCtx, "🚀 Running application...")
	err = a.Run(appCtx)
	if err != nil {
		logger.Error(appCtx, "❌ App.Run() returned error", zap.Error(err))
		return
	}

	logger.Info(appCtx, "👋 Application exited normally")
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
