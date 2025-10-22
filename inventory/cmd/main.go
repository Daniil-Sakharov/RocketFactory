package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/app"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func main() {
	// Перехватываем панику для отладки
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("🔥 PANIC: %v\n", r) //nolint:forbidigo // Early panic handler before logger init
			panic(r)                       // Повторно бросаем панику
		}
	}()

	// .env файл опционален:
	// - В локальной разработке: загружается из корня проекта или указанного пути
	// - В Docker: переменные передаются через environment (-e флаги)
	err := config.Load()
	if err != nil {
		fmt.Printf("❌ Failed to load config: %v\n", err) //nolint:forbidigo // Early config error before logger init
		panic(fmt.Errorf("error to load config: %w", err))
	}
	fmt.Println("✅ Config loaded") //nolint:forbidigo // Early success message before logger init

	appCtx, appCancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer appCancel()
	defer gracefulShutdown()

	closer.Configure(syscall.SIGINT, syscall.SIGTERM)

	fmt.Println("🏗️ Creating application...") //nolint:forbidigo // Early app creation message
	a, err := app.New(appCtx)
	if err != nil {
		fmt.Printf("❌ Failed to create app: %v\n", err) //nolint:forbidigo // Early app error
		logger.Error(appCtx, "❌ Не удалось создать приложение", zap.Error(err))
		return
	}
	fmt.Println("✅ Application created") //nolint:forbidigo // Early success message

	fmt.Println("🚀 Running application...") //nolint:forbidigo // Early run message
	err = a.Run(appCtx)
	if err != nil {
		fmt.Printf("❌ App.Run() returned error: %v\n", err) //nolint:forbidigo // Early run error
		logger.Error(appCtx, "❌ Ошибка при работе приложения", zap.Error(err))
		return
	}

	fmt.Println("👋 Application exited normally") //nolint:forbidigo // Early exit message
}

func gracefulShutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := closer.CloseAll(ctx); err != nil {
		logger.Error(ctx, "❌ Ошибка при завершении работы", zap.Error(err))
	}
}
