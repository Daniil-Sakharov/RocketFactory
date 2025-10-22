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

func NewRepository(_ context.Context, db *mongo.Database) *repository {
	collection := db.Collection("parts")

	// Create indexes synchronously but with error recovery
	// We use background context to avoid cancellation issues
	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "category", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	indexCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create indexes - if they already exist, MongoDB will silently skip
	_, _ = collection.Indexes().CreateMany(indexCtx, indexModel)
	// We ignore errors as indexes might already exist or creation can be retried later

	return &repository{collection: collection}
}
