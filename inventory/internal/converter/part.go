package converter

import (
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

// PartToProto конвертирует domain модель Part в protobuf модель
func PartToProto(part *model.Part) *inventoryv1.Part {
	if part == nil {
		return nil
	}

	protoPart := &inventoryv1.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToProto(part.Category),
		Dimensions:    DimensionsToProto(part.Dimensions),
		Manufacturer:  ManufacturerToProto(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      MetadataToProto(part.Metadata),
	}

	// Конвертируем timestamps
	if part.CreatedAt != nil {
		protoPart.CreatedAt = timestamppb.New(*part.CreatedAt)
	}
	if part.UpdatedAt != nil {
		protoPart.UpdatedAt = timestamppb.New(*part.UpdatedAt)
	}

	return protoPart
}

// PartsToProto конвертирует слайс domain моделей Part в слайс protobuf моделей
func PartsToProto(parts []*model.Part) []*inventoryv1.Part {
	if parts == nil {
		return nil
	}

	protoParts := make([]*inventoryv1.Part, 0, len(parts))
	for _, part := range parts {
		protoParts = append(protoParts, PartToProto(part))
	}

	return protoParts
}

// PartFromProto конвертирует protobuf модель Part в domain модель
func PartFromProto(protoPart *inventoryv1.Part) *model.Part {
	if protoPart == nil {
		return nil
	}

	part := &model.Part{
		Uuid:          protoPart.GetUuid(),
		Name:          protoPart.GetName(),
		Description:   protoPart.GetDescription(),
		Price:         protoPart.GetPrice(),
		StockQuantity: protoPart.GetStockQuantity(),
		Category:      CategoryFromProto(protoPart.GetCategory()),
		Dimensions:    DimensionsFromProto(protoPart.GetDimensions()),
		Manufacturer:  ManufacturerFromProto(protoPart.GetManufacturer()),
		Tags:          protoPart.GetTags(),
		Metadata:      MetadataFromProto(protoPart.GetMetadata()),
	}

	// Конвертируем timestamps
	if protoPart.GetCreatedAt() != nil {
		createdAt := protoPart.GetCreatedAt().AsTime()
		part.CreatedAt = &createdAt
	}
	if protoPart.GetUpdatedAt() != nil {
		updatedAt := protoPart.GetUpdatedAt().AsTime()
		part.UpdatedAt = &updatedAt
	}

	return part
}

// CategoryToProto конвертирует domain Category в protobuf Category
func CategoryToProto(category model.Category) inventoryv1.Category {
	switch category {
	case model.CATEGORY_ENGINE:
		return inventoryv1.Category_CATEGORY_ENGINE
	case model.CATEGORY_FUEL:
		return inventoryv1.Category_CATEGORY_FUEL
	case model.CATEGORY_PORTHOLE:
		return inventoryv1.Category_CATEGORY_PORTHOLE
	case model.CATEGORY_WING:
		return inventoryv1.Category_CATEGORY_WING
	default:
		return inventoryv1.Category_CATEGORY_UNSPECIFIED
	}
}

// CategoryFromProto конвертирует protobuf Category в domain Category
func CategoryFromProto(protoCategory inventoryv1.Category) model.Category {
	switch protoCategory {
	case inventoryv1.Category_CATEGORY_ENGINE:
		return model.CATEGORY_ENGINE
	case inventoryv1.Category_CATEGORY_FUEL:
		return model.CATEGORY_FUEL
	case inventoryv1.Category_CATEGORY_PORTHOLE:
		return model.CATEGORY_PORTHOLE
	case inventoryv1.Category_CATEGORY_WING:
		return model.CATEGORY_WING
	default:
		return model.CATEGORY_UNSPECIFIED
	}
}

// CategoriesToProto конвертирует слайс domain Categories в слайс protobuf Categories
func CategoriesToProto(categories []model.Category) []inventoryv1.Category {
	if categories == nil {
		return nil
	}

	protoCategories := make([]inventoryv1.Category, 0, len(categories))
	for _, category := range categories {
		protoCategories = append(protoCategories, CategoryToProto(category))
	}

	return protoCategories
}

// CategoriesFromProto конвертирует слайс protobuf Categories в слайс domain Categories
func CategoriesFromProto(protoCategories []inventoryv1.Category) []model.Category {
	if protoCategories == nil {
		return nil
	}

	categories := make([]model.Category, 0, len(protoCategories))
	for _, protoCategory := range protoCategories {
		categories = append(categories, CategoryFromProto(protoCategory))
	}

	return categories
}

// DimensionsToProto конвертирует domain Dimensions в protobuf Dimensions
func DimensionsToProto(dimensions *model.Dimensions) *inventoryv1.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &inventoryv1.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

// DimensionsFromProto конвертирует protobuf Dimensions в domain Dimensions
func DimensionsFromProto(protoDimensions *inventoryv1.Dimensions) *model.Dimensions {
	if protoDimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: protoDimensions.GetLength(),
		Width:  protoDimensions.GetWidth(),
		Height: protoDimensions.GetHeight(),
		Weight: protoDimensions.GetWeight(),
	}
}

// ManufacturerToProto конвертирует domain Manufacturer в protobuf Manufacturer
func ManufacturerToProto(manufacturer *model.Manufacturer) *inventoryv1.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &inventoryv1.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

// ManufacturerFromProto конвертирует protobuf Manufacturer в domain Manufacturer
func ManufacturerFromProto(protoManufacturer *inventoryv1.Manufacturer) *model.Manufacturer {
	if protoManufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    protoManufacturer.GetName(),
		Country: protoManufacturer.GetCountry(),
		Website: protoManufacturer.GetWebsite(),
	}
}

// MetadataToProto конвертирует domain metadata (map[string]interface{}) в protobuf metadata (map[string]*Value)
func MetadataToProto(metadata map[string]interface{}) map[string]*inventoryv1.Value {
	if metadata == nil {
		return nil
	}

	protoMetadata := make(map[string]*inventoryv1.Value, len(metadata))
	for key, value := range metadata {
		protoMetadata[key] = ValueToProto(value)
	}

	return protoMetadata
}

// MetadataFromProto конвертирует protobuf metadata (map[string]*Value) в domain metadata (map[string]interface{})
func MetadataFromProto(protoMetadata map[string]*inventoryv1.Value) map[string]interface{} {
	if protoMetadata == nil {
		return nil
	}

	metadata := make(map[string]interface{}, len(protoMetadata))
	for key, protoValue := range protoMetadata {
		metadata[key] = ValueFromProto(protoValue)
	}

	return metadata
}

// ValueToProto конвертирует interface{} в protobuf Value
func ValueToProto(value interface{}) *inventoryv1.Value {
	if value == nil {
		return nil
	}

	switch v := value.(type) {
	case string:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_StringValue{StringValue: v},
		}
	case int:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_Int64Value{Int64Value: int64(v)},
		}
	case int32:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_Int64Value{Int64Value: int64(v)},
		}
	case int64:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_Int64Value{Int64Value: v},
		}
	case float32:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_DoubleValue{DoubleValue: float64(v)},
		}
	case float64:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_DoubleValue{DoubleValue: v},
		}
	case bool:
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_BoolValue{BoolValue: v},
		}
	default:
		// Для неизвестных типов конвертируем в строку
		return &inventoryv1.Value{
			Value: &inventoryv1.Value_StringValue{StringValue: "unknown"},
		}
	}
}

// ValueFromProto конвертирует protobuf Value в interface{}
func ValueFromProto(protoValue *inventoryv1.Value) interface{} {
	if protoValue == nil {
		return nil
	}

	switch v := protoValue.GetValue().(type) {
	case *inventoryv1.Value_StringValue:
		return v.StringValue
	case *inventoryv1.Value_Int64Value:
		return v.Int64Value
	case *inventoryv1.Value_DoubleValue:
		return v.DoubleValue
	case *inventoryv1.Value_BoolValue:
		return v.BoolValue
	default:
		return nil
	}
}

// FilterToProto конвертирует domain PartsFilter в protobuf PartsFilter
func FilterToProto(filter *model.PartsFilter) *inventoryv1.PartsFilter {
	if filter == nil {
		return nil
	}

	return &inventoryv1.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToProto(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

// FilterFromProto конвертирует protobuf PartsFilter в domain PartsFilter
func FilterFromProto(protoFilter *inventoryv1.PartsFilter) *model.PartsFilter {
	if protoFilter == nil {
		return nil
	}

	return &model.PartsFilter{
		Uuids:                 protoFilter.GetUuids(),
		Names:                 protoFilter.GetNames(),
		Categories:            CategoriesFromProto(protoFilter.GetCategories()),
		ManufacturerCountries: protoFilter.GetManufacturerCountries(),
		Tags:                  protoFilter.GetTags(),
	}
}
