package part

import (
	"context"
	"errors"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (r *repository) ListParts(ctx context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	logger.Info(ctx, "🔍 ListParts called", zap.Any("filter", filter))

	mongoFilter := bson.M{}

	// Handle nil filter
	if filter == nil {
		logger.Info(ctx, "Filter is nil, returning all documents")
		filter = &model.PartsFilter{} // Create empty filter
	}

	if len(filter.Uuids) > 0 {
		mongoFilter["uuid"] = bson.M{"$in": filter.Uuids}
	}
	if len(filter.Names) > 0 {
		mongoFilter["name"] = bson.M{"$in": filter.Names}
	}
	if len(filter.Categories) > 0 {
		// Convert domain categories to repository format (strings)
		repoCategories := converter.CategoriesToRepo(filter.Categories)
		mongoFilter["category"] = bson.M{"$in": repoCategories}
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
		logger.Error(ctx, "Failed to decode cursor", zap.Error(err))
		return nil, fmt.Errorf("failed to parse: %w", err)
	}

	// Debug logging for tests
	logger.Info(ctx, "📊 ListParts result",
		zap.Int("found_count", len(repoParts)),
		zap.String("collection", r.collection.Name()),
		zap.Any("mongo_filter", mongoFilter))
	log.Printf("DEBUG: ListParts found %d parts in database, collection: %s\n", len(repoParts), r.collection.Name()) //nolint:forbidigo // Debug logging

	modelParts := converter.PartsToModel(repoParts)

	return modelParts, nil
}
