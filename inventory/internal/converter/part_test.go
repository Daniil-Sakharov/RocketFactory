package converter

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/Daniil-Sakharov/RocketFactory/inventory/internal/model"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

func TestPartToProto(t *testing.T) {
	now := time.Now()

	domainPart := &model.Part{
		Uuid:          "test-uuid",
		Name:          "Test Engine",
		Description:   "Test Description",
		Price:         1000000.0,
		StockQuantity: 5,
		Category:      model.CATEGORY_ENGINE,
		Dimensions: &model.Dimensions{
			Length: 100.0,
			Width:  50.0,
			Height: 50.0,
			Weight: 500.0,
		},
		Manufacturer: &model.Manufacturer{
			Name:    "TestCorp",
			Country: "USA",
			Website: "test.com",
		},
		Tags: []string{"test", "engine"},
		Metadata: map[string]interface{}{
			"power":  1000.0,
			"tested": true,
			"notes":  "test notes",
		},
		CreatedAt: &now,
		UpdatedAt: &now,
	}

	protoPart := PartToProto(domainPart)

	assert.NotNil(t, protoPart)
	assert.Equal(t, domainPart.Uuid, protoPart.Uuid)
	assert.Equal(t, domainPart.Name, protoPart.Name)
	assert.Equal(t, domainPart.Description, protoPart.Description)
	assert.Equal(t, domainPart.Price, protoPart.Price)
	assert.Equal(t, domainPart.StockQuantity, protoPart.StockQuantity)
	assert.Equal(t, inventoryv1.Category_CATEGORY_ENGINE, protoPart.Category)
	assert.NotNil(t, protoPart.Dimensions)
	assert.NotNil(t, protoPart.Manufacturer)
	assert.Equal(t, len(domainPart.Tags), len(protoPart.Tags))
	assert.NotNil(t, protoPart.Metadata)
	assert.NotNil(t, protoPart.CreatedAt)
	assert.NotNil(t, protoPart.UpdatedAt)
}

func TestPartToProto_Nil(t *testing.T) {
	protoPart := PartToProto(nil)
	assert.Nil(t, protoPart)
}

func TestPartFromProto(t *testing.T) {
	now := timestamppb.Now()

	protoPart := &inventoryv1.Part{
		Uuid:          "test-uuid",
		Name:          "Test Engine",
		Description:   "Test Description",
		Price:         1000000.0,
		StockQuantity: 5,
		Category:      inventoryv1.Category_CATEGORY_ENGINE,
		Dimensions: &inventoryv1.Dimensions{
			Length: 100.0,
			Width:  50.0,
			Height: 50.0,
			Weight: 500.0,
		},
		Manufacturer: &inventoryv1.Manufacturer{
			Name:    "TestCorp",
			Country: "USA",
			Website: "test.com",
		},
		Tags: []string{"test", "engine"},
		Metadata: map[string]*inventoryv1.Value{
			"power": {
				Value: &inventoryv1.Value_DoubleValue{DoubleValue: 1000.0},
			},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	domainPart := PartFromProto(protoPart)

	assert.NotNil(t, domainPart)
	assert.Equal(t, protoPart.Uuid, domainPart.Uuid)
	assert.Equal(t, protoPart.Name, domainPart.Name)
	assert.Equal(t, protoPart.Description, domainPart.Description)
	assert.Equal(t, protoPart.Price, domainPart.Price)
	assert.Equal(t, protoPart.StockQuantity, domainPart.StockQuantity)
	assert.Equal(t, model.CATEGORY_ENGINE, domainPart.Category)
	assert.NotNil(t, domainPart.Dimensions)
	assert.NotNil(t, domainPart.Manufacturer)
	assert.Equal(t, len(protoPart.Tags), len(domainPart.Tags))
	assert.NotNil(t, domainPart.Metadata)
	assert.NotNil(t, domainPart.CreatedAt)
	assert.NotNil(t, domainPart.UpdatedAt)
}

func TestPartFromProto_Nil(t *testing.T) {
	domainPart := PartFromProto(nil)
	assert.Nil(t, domainPart)
}

func TestPartsToProto(t *testing.T) {
	domainParts := []*model.Part{
		{
			Uuid:     "uuid-1",
			Name:     "Part 1",
			Category: model.CATEGORY_ENGINE,
		},
		{
			Uuid:     "uuid-2",
			Name:     "Part 2",
			Category: model.CATEGORY_FUEL,
		},
	}

	protoParts := PartsToProto(domainParts)

	assert.NotNil(t, protoParts)
	assert.Equal(t, len(domainParts), len(protoParts))
	assert.Equal(t, domainParts[0].Uuid, protoParts[0].Uuid)
	assert.Equal(t, domainParts[1].Uuid, protoParts[1].Uuid)
}

func TestPartsToProto_Nil(t *testing.T) {
	protoParts := PartsToProto(nil)
	assert.Nil(t, protoParts)
}

func TestCategoryToProto(t *testing.T) {
	tests := []struct {
		name     string
		category model.Category
		expected inventoryv1.Category
	}{
		{
			name:     "ENGINE",
			category: model.CATEGORY_ENGINE,
			expected: inventoryv1.Category_CATEGORY_ENGINE,
		},
		{
			name:     "FUEL",
			category: model.CATEGORY_FUEL,
			expected: inventoryv1.Category_CATEGORY_FUEL,
		},
		{
			name:     "PORTHOLE",
			category: model.CATEGORY_PORTHOLE,
			expected: inventoryv1.Category_CATEGORY_PORTHOLE,
		},
		{
			name:     "WING",
			category: model.CATEGORY_WING,
			expected: inventoryv1.Category_CATEGORY_WING,
		},
		{
			name:     "UNSPECIFIED",
			category: model.CATEGORY_UNSPECIFIED,
			expected: inventoryv1.Category_CATEGORY_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategoryToProto(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCategoryFromProto(t *testing.T) {
	tests := []struct {
		name     string
		category inventoryv1.Category
		expected model.Category
	}{
		{
			name:     "ENGINE",
			category: inventoryv1.Category_CATEGORY_ENGINE,
			expected: model.CATEGORY_ENGINE,
		},
		{
			name:     "FUEL",
			category: inventoryv1.Category_CATEGORY_FUEL,
			expected: model.CATEGORY_FUEL,
		},
		{
			name:     "PORTHOLE",
			category: inventoryv1.Category_CATEGORY_PORTHOLE,
			expected: model.CATEGORY_PORTHOLE,
		},
		{
			name:     "WING",
			category: inventoryv1.Category_CATEGORY_WING,
			expected: model.CATEGORY_WING,
		},
		{
			name:     "UNSPECIFIED",
			category: inventoryv1.Category_CATEGORY_UNSPECIFIED,
			expected: model.CATEGORY_UNSPECIFIED,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CategoryFromProto(tt.category)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFilterToProto(t *testing.T) {
	domainFilter := &model.PartsFilter{
		Uuids:                 []string{"uuid-1", "uuid-2"},
		Names:                 []string{"Engine", "Fuel"},
		Categories:            []model.Category{model.CATEGORY_ENGINE, model.CATEGORY_FUEL},
		ManufacturerCountries: []string{"USA", "Russia"},
		Tags:                  []string{"tag1", "tag2"},
	}

	protoFilter := FilterToProto(domainFilter)

	assert.NotNil(t, protoFilter)
	assert.Equal(t, domainFilter.Uuids, protoFilter.Uuids)
	assert.Equal(t, domainFilter.Names, protoFilter.Names)
	assert.Equal(t, len(domainFilter.Categories), len(protoFilter.Categories))
	assert.Equal(t, domainFilter.ManufacturerCountries, protoFilter.ManufacturerCountries)
	assert.Equal(t, domainFilter.Tags, protoFilter.Tags)
}

func TestFilterToProto_Nil(t *testing.T) {
	protoFilter := FilterToProto(nil)
	assert.Nil(t, protoFilter)
}

func TestFilterFromProto(t *testing.T) {
	protoFilter := &inventoryv1.PartsFilter{
		Uuids:                 []string{"uuid-1", "uuid-2"},
		Names:                 []string{"Engine", "Fuel"},
		Categories:            []inventoryv1.Category{inventoryv1.Category_CATEGORY_ENGINE, inventoryv1.Category_CATEGORY_FUEL},
		ManufacturerCountries: []string{"USA", "Russia"},
		Tags:                  []string{"tag1", "tag2"},
	}

	domainFilter := FilterFromProto(protoFilter)

	assert.NotNil(t, domainFilter)
	assert.Equal(t, protoFilter.Uuids, domainFilter.Uuids)
	assert.Equal(t, protoFilter.Names, domainFilter.Names)
	assert.Equal(t, len(protoFilter.Categories), len(domainFilter.Categories))
	assert.Equal(t, protoFilter.ManufacturerCountries, domainFilter.ManufacturerCountries)
	assert.Equal(t, protoFilter.Tags, domainFilter.Tags)
}

func TestFilterFromProto_Nil(t *testing.T) {
	domainFilter := FilterFromProto(nil)
	assert.Nil(t, domainFilter)
}

func TestValueToProto(t *testing.T) {
	tests := []struct {
		name     string
		value    interface{}
		expected interface{}
	}{
		{
			name:     "string",
			value:    "test",
			expected: &inventoryv1.Value_StringValue{StringValue: "test"},
		},
		{
			name:     "int",
			value:    42,
			expected: &inventoryv1.Value_Int64Value{Int64Value: 42},
		},
		{
			name:     "int64",
			value:    int64(42),
			expected: &inventoryv1.Value_Int64Value{Int64Value: 42},
		},
		{
			name:     "float64",
			value:    3.14,
			expected: &inventoryv1.Value_DoubleValue{DoubleValue: 3.14},
		},
		{
			name:     "bool",
			value:    true,
			expected: &inventoryv1.Value_BoolValue{BoolValue: true},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValueToProto(tt.value)
			assert.NotNil(t, result)
			assert.Equal(t, tt.expected, result.Value)
		})
	}
}

func TestValueToProto_Nil(t *testing.T) {
	result := ValueToProto(nil)
	assert.Nil(t, result)
}

func TestValueFromProto(t *testing.T) {
	tests := []struct {
		name     string
		value    *inventoryv1.Value
		expected interface{}
	}{
		{
			name: "string",
			value: &inventoryv1.Value{
				Value: &inventoryv1.Value_StringValue{StringValue: "test"},
			},
			expected: "test",
		},
		{
			name: "int64",
			value: &inventoryv1.Value{
				Value: &inventoryv1.Value_Int64Value{Int64Value: 42},
			},
			expected: int64(42),
		},
		{
			name: "double",
			value: &inventoryv1.Value{
				Value: &inventoryv1.Value_DoubleValue{DoubleValue: 3.14},
			},
			expected: 3.14,
		},
		{
			name: "bool",
			value: &inventoryv1.Value{
				Value: &inventoryv1.Value_BoolValue{BoolValue: true},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValueFromProto(tt.value)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestValueFromProto_Nil(t *testing.T) {
	result := ValueFromProto(nil)
	assert.Nil(t, result)
}
