package mongo

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func connectMongoClient(ctx context.Context, uri string) (*mongo.Client, error) {
	// MongoDB в Docker может инициализироваться до 30 секунд
	maxRetries := 15
	retryDelay := 2 * time.Second

	var client *mongo.Client
	var err error

	for attempt := 1; attempt <= maxRetries; attempt++ {
		client, err = mongo.Connect(ctx, options.Client().ApplyURI(uri))
		if err != nil {
			if attempt < maxRetries {
				timer := time.NewTimer(retryDelay)
				select {
				case <-timer.C:
				case <-ctx.Done():
					timer.Stop()
					return nil, errors.Wrap(ctx.Err(), "context cancelled while connecting to mongo")
				}
				continue
			}
			return nil, errors.Errorf("failed to connect to mongo after %d attempts: %v", maxRetries, err)
		}

		// Проверяем подключение
		err = client.Ping(ctx, readpref.Primary())
		if err == nil {
			return client, nil
		}

		// Если ping не прошел, закрываем клиента и пробуем снова
		// Игнорируем ошибку disconnect - попытаемся переподключиться
		_ = client.Disconnect(ctx) //nolint:gosec // Error is intentionally ignored

		if attempt < maxRetries {
			timer := time.NewTimer(retryDelay)
			select {
			case <-timer.C:
			case <-ctx.Done():
				timer.Stop()
				return nil, errors.Wrap(ctx.Err(), "context cancelled while connecting to mongo")
			}
		}
	}

	return nil, errors.Errorf("failed to ping mongo after %d attempts: %v", maxRetries, err)
}
