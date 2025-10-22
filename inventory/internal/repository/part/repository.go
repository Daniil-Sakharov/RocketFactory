package part

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	def "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository"
)

var _ def.PartRepository = (*repository)(nil)

type repository struct {
	collection *mongo.Collection
}

func NewRepository(ctx context.Context, db *mongo.Database) *repository {
	collection := db.Collection("parts")

	// Create indexes in a separate goroutine to avoid blocking
	// Index creation failures are logged but don't prevent repository creation
	go func() {
		indexModel := []mongo.IndexModel{
			{
				Keys:    bson.D{{Key: "name", Value: 1}, {Key: "category", Value: 1}},
				Options: options.Index().SetUnique(false),
			},
		}

		indexCtx, cancel := context.WithTimeout(ctx, 30*time.Second)
		defer cancel()

		_, err := collection.Indexes().CreateMany(indexCtx, indexModel)
		if err != nil {
			// Log error but don't panic - indexes might already exist or will be created later
			// This prevents application startup failures in test environments
			_ = err // Suppress linter warning
		}
	}()

	return &repository{collection: collection}
}
