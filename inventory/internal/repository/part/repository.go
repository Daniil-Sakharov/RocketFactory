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

func NewRepository(db *mongo.Database) *repository {
	collection := db.Collection("parts")

	indexModel := []mongo.IndexModel{
		{
			Keys:    bson.D{{Key: "name", Value: 1}, {Key: "category", Value: 1}},
			Options: options.Index().SetUnique(false),
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateMany(ctx, indexModel)
	if err != nil {
		panic(err)
	}

	return &repository{collection: collection}
}
