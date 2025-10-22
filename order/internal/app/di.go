package app

import (
	"context"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	orderAPI "github.com/Daniil-Sakharov/RocketFactory/order/internal/api/order/v1"
	grpcClient "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc"
	inventoryClient "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc/payment/v1"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/config"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/repository"
	orderRepo "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/order"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/service"
	orderService "github.com/Daniil-Sakharov/RocketFactory/order/internal/service/order"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator/pg"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	inventoryClient grpcClient.InventoryClient
	paymentClient   grpcClient.PaymentClient
	orderService    service.OrderService
	orderRepository repository.OrderRepository
	postgresDB      *sqlx.DB
	orderV1API      orderV1.Handler
	migrator        migrator.Migrator
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) OrderV1API(ctx context.Context) orderV1.Handler {
	if d.orderV1API == nil {
		d.orderV1API = orderAPI.NewAPI(d.OrderService(ctx))
	}
	return d.orderV1API
}

func (d *diContainer) OrderService(ctx context.Context) service.OrderService {
	if d.orderService == nil {
		d.orderService = orderService.NewService(
			d.OrderRepository(ctx),
			d.InventoryClient(),
			d.PaymentClient(),
		)
	}
	return d.orderService
}

func (d *diContainer) PaymentClient() grpcClient.PaymentClient {
	if d.paymentClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().PaymentGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("Ошибка в подключении к Payment Service: %s\n", err.Error()))
		}
		paymentGRPCStub := paymentV1.NewPaymentServiceClient(conn)
		closer.AddNamed("PaymentClient", func(ctx context.Context) error {
			return conn.Close()
		})
		d.paymentClient = paymentClient.NewClient(paymentGRPCStub)
	}
	return d.paymentClient
}

func (d *diContainer) InventoryClient() grpcClient.InventoryClient {
	if d.inventoryClient == nil {
		conn, err := grpc.NewClient(config.AppConfig().InventoryGRPC.Address(), grpc.WithTransportCredentials(insecure.NewCredentials()))
		if err != nil {
			panic(fmt.Sprintf("Ошибка в подключении к Inventory Service: %s\n", err.Error()))
		}
		inventoryGRPCStub := inventoryV1.NewInventoryServiceClient(conn)
		closer.AddNamed("InventoryClient", func(ctx context.Context) error {
			return conn.Close()
		})
		d.inventoryClient = inventoryClient.NewClient(inventoryGRPCStub)
	}
	return d.inventoryClient
}

func (d *diContainer) Migrator(ctx context.Context) migrator.Migrator {
	if d.migrator == nil {
		db := d.PostgresDB(ctx)
		d.migrator = pg.NewMigrator(db.DB, config.AppConfig().PostgresDB.MigrationsDir())
	}
	return d.migrator
}

func (d *diContainer) OrderRepository(ctx context.Context) repository.OrderRepository {
	if d.orderRepository == nil {
		d.orderRepository = orderRepo.NewRepository(d.PostgresDB(ctx))
	}
	return d.orderRepository
}

func (d *diContainer) PostgresDB(ctx context.Context) *sqlx.DB {
	if d.postgresDB == nil {
		db, err := sqlx.Connect("pgx", config.AppConfig().PostgresDB.URI())
		if err != nil {
			panic(fmt.Sprintf("Ошибка в подключении к PostgreSQL: %s\n", err.Error()))
		}

		err = db.Ping()
		if err != nil {
			panic(fmt.Sprintf("Ошибка в соединении с PostgreSQL: %s\n", err.Error()))
		}

		closer.AddNamed("PostgreSQL", func(ctx context.Context) error {
			return db.Close()
		})

		d.postgresDB = db
	}
	return d.postgresDB
}
