package app

import (
	"context"
	"fmt"

	"github.com/IBM/sarama"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"

	httpClient "github.com/Daniil-Sakharov/RocketFactory/notification/internal/client/http"
	telegramClient "github.com/Daniil-Sakharov/RocketFactory/notification/internal/client/http/telegram"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/config"
	kafkaConverter "github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/converter/kafka/decoder"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service"
	orderPaidConsumer "github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/consumer/order_paid_consumer"
	shipAssemledConsumer "github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/consumer/ship_assembly_consumer"
	"github.com/Daniil-Sakharov/RocketFactory/notification/internal/service/telegram"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/closer"
	wrappedKafka "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka"
	wrappedKafkaConsumer "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/kafka/consumer"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
	kafkaMiddleware "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/middleware/kafka"
)

type diContainer struct {
	telegramService             service.TelegramService
	botService                  service.BotService
	orderPaidConsumerService    service.OrderPaidConsumerService
	shipAssemblyConsumerService service.ShipAssemblyConsumerService

	orderPaidDecoder     kafkaConverter.OrderDecoder
	shipAssembledDecoder kafkaConverter.AssemblyDecoder

	orderPaidConsumerGroup     sarama.ConsumerGroup
	shipAssembledConsumerGroup sarama.ConsumerGroup
	orderPaidConsumer          wrappedKafka.Consumer
	shipAssembledConsumer      wrappedKafka.Consumer

	telegramBot    *bot.Bot
	telegramClient httpClient.TelegramClient
	templateEngine *telegram.TemplateEngine
}

func NewDiContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) TelegramBot() *bot.Bot {
	if d.telegramBot == nil {
		opts := []bot.Option{
			bot.WithDefaultHandler(func(ctx context.Context, b *bot.Bot, update *models.Update) {
				logger.Info(ctx, "Received unhandled update")
			}),
		}

		b, err := bot.New(config.AppConfig().TelegramBot.Token(), opts...)
		if err != nil {
			panic(fmt.Sprintf("failed to create telegram bot: %s", err.Error()))
		}

		d.telegramBot = b
	}
	return d.telegramBot
}

func (d *diContainer) TelegramClient() httpClient.TelegramClient {
	if d.telegramClient == nil {
		d.telegramClient = telegramClient.NewClient(d.TelegramBot())
	}
	return d.telegramClient
}

func (d *diContainer) TemplateEngine() *telegram.TemplateEngine {
	if d.templateEngine == nil {
		engine, err := telegram.NewTemplateEngine()
		if err != nil {
			panic(fmt.Sprintf("failed to create template engine: %v", err))
		}
		d.templateEngine = engine
	}
	return d.templateEngine
}

func (d *diContainer) TelegramService() service.TelegramService {
	if d.telegramService == nil {
		d.telegramService = telegram.NewService(
			d.TelegramClient(),
			d.TemplateEngine(),
		)
	}
	return d.telegramService
}

func (d *diContainer) BotService() service.BotService {
	if d.botService == nil {
		d.botService = telegram.NewBotService(
			d.TelegramBot(),
			d.TemplateEngine(),
		)
	}
	return d.botService
}

func (d *diContainer) OrderPaidDecoder() kafkaConverter.OrderDecoder {
	if d.orderPaidDecoder == nil {
		d.orderPaidDecoder = decoder.NewOrderDecoder()
	}
	return d.orderPaidDecoder
}

func (d *diContainer) ShipAssembledDecoder() kafkaConverter.AssemblyDecoder {
	if d.shipAssembledDecoder == nil {
		d.shipAssembledDecoder = decoder.NewAssemblyDecoder()
	}
	return d.shipAssembledDecoder
}

func (d *diContainer) OrderPaidConsumerGroup() sarama.ConsumerGroup {
	if d.orderPaidConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().OrderConsumer.GroupID(),
			config.AppConfig().OrderConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create order_paid consumer group: %s", err.Error()))
		}

		closer.AddNamed("Kafka OrderPaid consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.orderPaidConsumerGroup = consumerGroup
	}
	return d.orderPaidConsumerGroup
}

func (d *diContainer) ShipAssembledConsumerGroup() sarama.ConsumerGroup {
	if d.shipAssembledConsumerGroup == nil {
		consumerGroup, err := sarama.NewConsumerGroup(
			config.AppConfig().Kafka.Brokers(),
			config.AppConfig().AssemblyConsumer.GroupID(),
			config.AppConfig().AssemblyConsumer.Config(),
		)
		if err != nil {
			panic(fmt.Sprintf("failed to create ship_assembled consumer group: %s", err.Error()))
		}

		closer.AddNamed("Kafka ShipAssembled consumer group", func(ctx context.Context) error {
			return consumerGroup.Close()
		})

		d.shipAssembledConsumerGroup = consumerGroup
	}
	return d.shipAssembledConsumerGroup
}

func (d *diContainer) OrderPaidConsumer() wrappedKafka.Consumer {
	if d.orderPaidConsumer == nil {
		d.orderPaidConsumer = wrappedKafkaConsumer.NewConsumer(
			d.OrderPaidConsumerGroup(),
			[]string{
				config.AppConfig().OrderConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.orderPaidConsumer
}

func (d *diContainer) ShipAssembledConsumer() wrappedKafka.Consumer {
	if d.shipAssembledConsumer == nil {
		d.shipAssembledConsumer = wrappedKafkaConsumer.NewConsumer(
			d.ShipAssembledConsumerGroup(),
			[]string{
				config.AppConfig().AssemblyConsumer.Topic(),
			},
			logger.Logger(),
			kafkaMiddleware.Logging(logger.Logger()),
		)
	}
	return d.shipAssembledConsumer
}

func (d *diContainer) OrderPaidConsumerService() service.OrderPaidConsumerService {
	if d.orderPaidConsumerService == nil {
		d.orderPaidConsumerService = orderPaidConsumer.NewService(
			d.OrderPaidConsumer(),
			d.OrderPaidDecoder(),
			d.TelegramService(),
		)
	}
	return d.orderPaidConsumerService
}

func (d *diContainer) ShipAssemblyConsumerService() service.ShipAssemblyConsumerService {
	if d.shipAssemblyConsumerService == nil {
		d.shipAssemblyConsumerService = shipAssemledConsumer.NewService(
			d.ShipAssembledConsumer(),
			d.ShipAssembledDecoder(),
			d.TelegramService(),
		)
	}
	return d.shipAssemblyConsumerService
}
