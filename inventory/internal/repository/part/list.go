package part

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	mongoFilter := bson.M{}

	if len(filter.Uuids) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.Uuids}
	}
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		mongoFilter["category"] = bson.M{"$in": filter.Categories}
	}
	if len(filter.ManufacturerCountries) > 0 {
		mongoFilter["manufacturer_country"] = bson.M{"$in": filter.ManufacturerCountries}
	}
	if len(filter.Tags) > 0 {
		mongoFilter["tags"] = bson.M{"$in": filter.Tags}
	}

	var repoParts []*repoModel.Part

	cursor, err := r.collection.Find(ctx, mongoFilter)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartsNotFound
		}
		return nil, fmt.Errorf("failed to find parts %w", err)
	}

	defer func() {
		err = cursor.Close(ctx)
		if err != nil {
			log.Println("failed to close cursor")
		}
	}()

	err = cursor.All(ctx, &repoParts)
	if err != nil {
		return nil, fmt.Errorf("failed to parse")
	}
	modelParts := converter.PartsToModel(repoParts)

	return modelParts, nil
}
