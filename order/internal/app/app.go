package app

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/http/health"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
)

type App struct {
	diContainer *diContainer
	httpServer  http.Server
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	errCh := make(chan error, 2)

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		if err := a.runHTTPServer(ctx); err != nil {
			errCh <- errors.Errorf("HHTP server crashed: %v", err)
		}
	}()
	go func() {
		if err := a.runConsumer(ctx); err != nil {
			errCh <- errors.Errorf("consumer crashed: %v", err)
		}
	}()

	select {
	case <-ctx.Done():
		logger.Info(ctx, "Shutdown signal received")
	case err := <-errCh:
		logger.Error(ctx, "Component crashed, shutting down", zap.Error(err))
		cancel()
		<-ctx.Done()
		return err
	}
	return nil
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initLogger,
		a.initDI,
		a.StartMigrations,
		a.initCloser,
		a.initHTTPServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *App) initDI(_ context.Context) error {
	a.diContainer = NewDiContainer()
	return nil
}

func (a *App) StartMigrations(ctx context.Context) error {
	migrator := a.diContainer.Migrator(ctx)
	if err := migrator.Up(ctx); err != nil {
		return err
	}
	return nil
}

func (a *App) initLogger(_ context.Context) error {
	return logger.Init(
		config.AppConfig().Logger.Level(),
		config.AppConfig().Logger.AsJson(),
	)
}

func (a *App) initCloser(_ context.Context) error {
	closer.SetLogger(logger.Logger())
	return nil
}

func (a *App) initHTTPServer(ctx context.Context) error {
	mux := http.NewServeMux()

	server, err := orderV1.NewServer(a.diContainer.OrderV1API(ctx))
	if err != nil {
		return err
	}

	mux.Handle("/health", health.NewHandler(health.Config{
		ServiceName: "order-service",
		Version:     "1.0.0",
	}))
	mux.Handle("/api/", server)

	a.httpServer = http.Server{
		Addr:         config.AppConfig().OrderHTTP.Address(),
		Handler:      mux,
		ReadTimeout:  config.AppConfig().OrderHTTP.ReadTimeout(),
		WriteTimeout: config.AppConfig().OrderHTTP.WriteTimeout(),
		IdleTimeout:  config.AppConfig().OrderHTTP.IdleTimeout(),
	}

	closer.AddNamed("HTTP server", func(ctx context.Context) error {
		return a.httpServer.Close()
	})

	return nil
}

func (a *App) runHTTPServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 HTTP OrderService server listening on %s", config.AppConfig().OrderHTTP.Address()))

	if err := a.httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (a *App) runConsumer(ctx context.Context) error {
	logger.Info(ctx, "🚀 OrderPaid Kafka consumer starting")

	err := a.diContainer.AssemblyConsumerService(ctx).RunConsumer(ctx)
	if err != nil {
		return err
	}

	return nil
}
