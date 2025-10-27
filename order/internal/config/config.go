package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Daniil-Sakharov/RocketFactory/order/internal/config/env"
)

var appConfig *config

type config struct {
	Logger           LoggerConfig
	OrderHTTP        OrderHTTPConfig
	InventoryGRPC    InventoryGRPCConfig
	PaymentGRPC      PaymentGRPCConfig
	PostgresDB       PostgresConfig
	Kafka            KafkaConfig
	AssemblyConsumer AssemblyConsumerConfig
	OrderProducer    OrderProducerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	inventoryGRPCCfg, err := env.NewInventoryGRPCConfig()
	if err != nil {
		return err
	}

	paymentGRPCCfg, err := env.NewPaymentGRPCConfig()
	if err != nil {
		return err
	}

	orderHHTPCfg, err := env.NewOrderHTTPConfig()
	if err != nil {
		return err
	}

	postgresCfg, err := env.NewPostgresConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	producerCfg, err := env.NewOrderProduceConfig()
	if err != nil {
		return err
	}

	consumerCfg, err := env.NewAssemblyConsumerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:           loggerCfg,
		OrderHTTP:        orderHHTPCfg,
		InventoryGRPC:    inventoryGRPCCfg,
		PaymentGRPC:      paymentGRPCCfg,
		PostgresDB:       postgresCfg,
		Kafka:            kafkaCfg,
		OrderProducer:    producerCfg,
		AssemblyConsumer: consumerCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
