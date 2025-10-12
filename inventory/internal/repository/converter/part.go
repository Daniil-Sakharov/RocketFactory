package converter

import (
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
)

func PartsToModel(repoParts []*repoModel.Part) []*model.Part {
	if repoParts == nil {
		return nil
	}

	result := make([]*model.Part, 0, len(repoParts))
	for _, repoPart := range repoParts {
		result = append(result, PartToModel(repoPart))
	}
	return result
}

func FilterToRepoModel(filter *model.PartsFilter) *repoModel.PartsFilter {
	if filter == nil {
		return nil
	}

	return &repoModel.PartsFilter{
		Uuids:                 filter.Uuids,
		Names:                 filter.Names,
		Categories:            CategoriesToRepo(filter.Categories),
		ManufacturerCountries: filter.ManufacturerCountries,
		Tags:                  filter.Tags,
	}
}

func CategoriesToModel(categories []repoModel.Category) []model.Category {
	if categories == nil {
		return nil
	}

	result := make([]model.Category, len(categories))
	for i, cat := range categories {
		result[i] = model.Category(cat)
	}
	return result
}

func CategoriesToRepo(categories []model.Category) []repoModel.Category {
	if categories == nil {
		return nil
	}

	result := make([]repoModel.Category, len(categories))
	for i, cat := range categories {
		result[i] = repoModel.Category(cat)
	}
	return result
}

func PartToModel(part *repoModel.Part) *model.Part {
	if part == nil {
		return nil
	}

	return &model.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToModel(part.Category),
		Dimensions:    DimensionsToModel(part.Dimensions),
		Manufacturer:  ManufacturerToModel(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ManufacturerToModel(manufacturer *repoModel.Manufacturer) *model.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &model.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func CategoryToModel(category repoModel.Category) model.Category {
	return model.Category(category)
}

func DimensionsToModel(dimensions *repoModel.Dimensions) *model.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &model.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}

func PartToRepoModel(part *model.Part) *repoModel.Part {
	if part == nil {
		return nil
	}

	return &repoModel.Part{
		Uuid:          part.Uuid,
		Name:          part.Name,
		Description:   part.Description,
		Price:         part.Price,
		StockQuantity: part.StockQuantity,
		Category:      CategoryToRepo(part.Category),
		Dimensions:    DimensionsToRepo(part.Dimensions),
		Manufacturer:  ManufacturerToRepo(part.Manufacturer),
		Tags:          part.Tags,
		Metadata:      part.Metadata,
		CreatedAt:     part.CreatedAt,
		UpdatedAt:     part.UpdatedAt,
	}
}

func ManufacturerToRepo(manufacturer *model.Manufacturer) *repoModel.Manufacturer {
	if manufacturer == nil {
		return nil
	}

	return &repoModel.Manufacturer{
		Name:    manufacturer.Name,
		Country: manufacturer.Country,
		Website: manufacturer.Website,
	}
}

func CategoryToRepo(category model.Category) repoModel.Category {
	return repoModel.Category(category)
}

func DimensionsToRepo(dimensions *model.Dimensions) *repoModel.Dimensions {
	if dimensions == nil {
		return nil
	}

	return &repoModel.Dimensions{
		Length: dimensions.Length,
		Width:  dimensions.Width,
		Height: dimensions.Height,
		Weight: dimensions.Weight,
	}
}
