package app

import (
    // Standard library
    "context"
    "errors"
    "fmt"
    "net"

    // Third-party
    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/reflection"

    // Module-local
    "github.com/Daniil-Sakharov/RocketFactory/payment/internal/config"
    "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
    "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/grpc/health"
    "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
    paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

type App struct {
	diContainer *diContainer
	grpcServer  *grpc.Server
	listener    net.Listener
}

func New(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.InitDeps(ctx)
	if err != nil {
        return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	return a.runGRPCServer(ctx)
}

func (a *App) InitDeps(ctx context.Context) error {
    inits := []func(context.Context) error{
		a.initDI,
		a.initLogger,
		a.initCloser,
		a.initListener,
		a.initGRPCServer,
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

func (a *App) initListener(_ context.Context) error {
	listener, err := net.Listen("tcp", config.AppConfig().Payment.Address())
	if err != nil {
		return err
	}
	closer.AddNamed("TCP listener", func(ctx context.Context) error {
		lerr := listener.Close()
		if lerr != nil && !errors.Is(lerr, net.ErrClosed) {
			return lerr
		}

		return nil
	})

	a.listener = listener

	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(grpc.Creds(insecure.NewCredentials()))
	closer.AddNamed("gRPC server", func(ctx context.Context) error {
		a.grpcServer.GracefulStop()
		return nil
	})

	reflection.Register(a.grpcServer)

	health.RegisterService(a.grpcServer)

	paymentv1.RegisterPaymentServiceServer(a.grpcServer, a.diContainer.PaymentAPI(ctx))

	return nil
}

func (a *App) runGRPCServer(ctx context.Context) error {
	logger.Info(ctx, fmt.Sprintf("🚀 gRPC PaymentService server listening on %s", config.AppConfig().Payment.Address()))

	err := a.grpcServer.Serve(a.listener)
	if err != nil {
		return err
	}

	return nil
}
