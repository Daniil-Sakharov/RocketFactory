package config

import (
	"os"

	"github.com/joho/godotenv"

	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/config/env"
)

var appConfig *config

type config struct {
	Logger           LoggerConfig
	Kafka            KafkaConfig
	TelegramBot      TelegramBotConfig
	OrderConsumer    OrderConsumerConfig
	AssemblyConsumer AssemblyConsumerConfig
}

func Load(path ...string) error {
	err := godotenv.Load(path...)
	if err != nil && !os.IsNotExist(err) {
		return err
	}

	assemblyCfg, err := env.NewAssemblyConsumerConfig()
	if err != nil {
		return err
	}
	orderCfg, err := env.NewOrderConsumerConfig()
	if err != nil {
		return err
	}
	tokenCfg, err := env.NewTelegramBotConfig()
	if err != nil {
		return err
	}

	kafkaCfg, err := env.NewKafkaConfig()
	if err != nil {
		return err
	}

	loggerCfg, err := env.NewLoggerConfig()
	if err != nil {
		return err
	}

	appConfig = &config{
		Logger:           loggerCfg,
		Kafka:            kafkaCfg,
		OrderConsumer:    orderCfg,
		AssemblyConsumer: assemblyCfg,
		TelegramBot:      tokenCfg,
	}

	return nil
}

func AppConfig() *config {
	return appConfig
}
