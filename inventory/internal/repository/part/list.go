package part

import (
	"context"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/converter"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/model"
)

func (r *repository) ListParts(_ context.Context, filter *model.PartsFilter) ([]*model.Part, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	repoFilter := converter.FilterToRepoModel(filter)

	repoParts := make([]*repoModel.Part, 0, len(r.data))
	for _, part := range r.data {
		partCopy := part
		if matchesFilter(&partCopy, repoFilter) {
			repoParts = append(repoParts, &partCopy)
		}
	}

	return converter.PartsToModel(repoParts), nil
}

func matchesFilter(part *repoModel.Part, filter *repoModel.PartsFilter) bool {
	if filter == nil {
		return true
	}

	if len(filter.Uuids) > 0 && !contains(filter.Uuids, part.Uuid) {
		return false
	}

	if len(filter.Names) > 0 && !contains(filter.Names, part.Name) {
		return false
	}

	if len(filter.Categories) > 0 && !containsCategory(filter.Categories, part.Category) {
		return false
	}

	if len(filter.ManufacturerCountries) > 0 && !contains(filter.ManufacturerCountries, part.Manufacturer.Country) {
		return false
	}
	if len(filter.Tags) > 0 && !hasAnyTag(filter.Tags, part.Tags) {
		return false
	}

	return true
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func containsCategory(slice []repoModel.Category, item repoModel.Category) bool {
	for _, c := range slice {
		if c == item {
			return true
		}
	}
	return false
}

func hasAnyTag(filterTags, partTags []string) bool {
	for _, filterTag := range filterTags {
		for _, partTag := range partTags {
			if filterTag == partTag {
				return true
			}
		}
	}
	return false
}
