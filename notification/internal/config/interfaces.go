package config

import "github.com/IBM/sarama"

type LoggerConfig interface {
	Level() string
	AsJson() bool
}

type TelegramBotConfig interface {
	Token() string
}

type KafkaConfig interface {
	Brokers() []string
}

type OrderConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
type AssemblyConsumerConfig interface {
	Topic() string
	GroupID() string
	Config() *sarama.Config
}
