package app

import (
	"context"
	"fmt"
	kafkaConverter "github.com/Daniil-Sakharov/RocketFactory/order/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/converter/kafka/decoder"
	wrappedKafka "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	kafkaMiddleware "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/middleware/kafka"
	"github.com/IBM/sarama"

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
	assemblyConsumer "github.com/Daniil-Sakharov/RocketFactory/order/internal/service/consumer/assembly_consumer"
	orderService "github.com/Daniil-Sakharov/RocketFactory/order/internal/service/order"
	orderProducer "github.com/Daniil-Sakharov/RocketFactory/order/internal/service/producer/order_producer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	wrappedKafkaConsumer "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka/consumer"
	wrappedKafkaProducer "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka/producer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/migrator/pg"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

type diContainer struct {
	inventoryClient         grpcClient.InventoryClient
	paymentClient           grpcClient.PaymentClient
	orderService            service.OrderService
	assemblyConsumerService service.AssemblyConsumerService
	orderProducerService    service.OrderProducerService
	orderRepository         repository.OrderRepository
	postgresDB              *sqlx.DB
	orderV1API              orderV1.Handler
	migrator                migrator.Migrator
	consumerGroup           sarama.ConsumerGroup
	assemblyConsumer        wrappedKafka.Consumer
	orderProducer           wrappedKafka.Producer
	assemblyDecoder         kafkaConverter.AssemblyDecoder
	syncProducer            sarama.SyncProducer
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
			d.OrderProducerService(),
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

func (d *diContainer) AssemblyConsumerService(ctx context.Context) service.AssemblyConsumerService {
	if d.assemblyConsumerService == nil {
		d.assemblyConsumerService = assemblyConsumer.NewService(
			d.AssemblyConsumer(),
			d.AssemblyDecoder(),
			d.OrderService(ctx),
			d.OrderRepository(ctx),
			)
	}
	return d.assemblyConsumerService
}

func (d *diContainer) AssemblyConsumer() wrappedKafka.Consumer {
	if d.assemblyConsumer == nil {
		d.assemblyConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ConsumerGroup(),
			[]string {
				config.AppConfig().OrderProducer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
			)
	}
	return d.assemblyConsumer
}

func (d *diContainer) ConsumerGroup() sarama.ConsumerGroup {
	if d.consumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyConsumer.GroupID(),
			config.AppConfig().AssemblyConsumer.Config(),
			)
		if err != nil {
			panic(fmt.Sprintf("failed to create consumer group: %s\n", err.Error()))
		}
		closer.AddNamed("Kafka consumer group", func(ctx context.Context) error {
			return d.consumerGroup.Close()
		})

		d.consumerGroup = consumerGroup
	}
	return d.consumerGroup
}

func (d *diContainer) AssemblyDecoder() kafkaConverter.AssemblyDecoder {
	if d.assemblyDecoder == nil {
		d.assemblyDecoder = decoder.NewAssemblyDecoder()
	}
	return d.assemblyDecoder
}

func (d *diContainer) OrderProducerService() service.OrderProducerService {
	if d.orderProducerService == nil {
		d.orderProducerService = orderProducer.NewService(d.OrderProducer())
	}
	return d.orderProducerService
}

func (d *diContainer) OrderProducer() wrappedKafka.Producer {
	if d.orderProducer == nil {
		d.orderProducer = wrappedKafkaProducer.NewProducer(
			d.SyncProducer(),
			config.AppConfig().OrderProducer.Topic(),
			logger.Logger(),
			)
	}
	return d.orderProducer
}

func (d *diContainer) SyncProducer() sarama.SyncProducer {
	if d.syncProducer == nil {
		p, err := sarama.NewSyncProducer(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderProducer.Config(),
			)
		if err != nil {
			panic(fmt.Sprintf("failed to create sync producer"))
		}
		closer.AddNamed("Kafka sync producer", func(ctx context.Context) error{
			return p.Close()
		})

		d.syncProducer = p
	}
	return d.syncProducer
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
