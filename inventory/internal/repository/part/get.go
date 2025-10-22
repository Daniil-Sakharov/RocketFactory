package part

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.uber.org/zap"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
	"github.com/Daniil-Sakharov/RocketFactory/platform/pkg/logger"
)

func (r *repository) GetPart(ctx context.Context, uuid string) (*model.Part, error) {
	logger.Info(ctx, "🔍 GetPart called", zap.String("uuid", uuid))

	var repoPart repoModel.Part
	err := r.collection.FindOne(ctx, bson.M{"uuid": uuid}).Decode(&repoPart)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, model.ErrPartNotFound
		}
		return nil, fmt.Errorf("failed to get part: %w", err)
	}

	modelPart := converter.PartToModel(&repoPart)
	return modelPart, nil
}
