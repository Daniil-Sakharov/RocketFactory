package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/timestamppb"

	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

type inventoryService struct {
	inventoryv1.UnimplementedInventoryServiceServer

	storage map[string]*inventoryv1.Part
	mu      *sync.RWMutex
}

func (i *inventoryService) GetPart(_ context.Context, req *inventoryv1.GetPartRequest) (*inventoryv1.GetPartResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	res, ok := i.storage[req.Uuid]
	if !ok {
		return nil, errors.New("NotFound")
	}

	return &inventoryv1.GetPartResponse{
		Part: res,
	}, nil
}

func (i *inventoryService) ListParts(_ context.Context, req *inventoryv1.ListPartsRequest) (*inventoryv1.ListPartsResponse, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()

	result := make([]*inventoryv1.Part, 0, len(i.storage))

	for _, part := range i.storage {
		if matchesFilter(part, req.Filter) {
			result = append(result, part)
		}
	}

	return &inventoryv1.ListPartsResponse{
		Parts: result,
	}, nil
}

func matchesFilter(part *inventoryv1.Part, filter *inventoryv1.PartsFilter) bool {
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

func containsCategory(slice []inventoryv1.Category, item inventoryv1.Category) bool {
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

func createTestData() map[string]*inventoryv1.Part {
	now := timestamppb.Now()

	return map[string]*inventoryv1.Part{
		"550e8400-e29b-41d4-a716-446655440001": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440001",
			Name:          "Ракетный двигатель RD-180",
			Description:   "Мощный жидкостный ракетный двигатель",
			Price:         15000000.0,
			StockQuantity: 3,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 350.0,
				Width:  240.0,
				Height: 240.0,
				Weight: 5480.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Энергомаш",
				Country: "Russia",
				Website: "www.npoenergomash.ru",
			},
			Tags: []string{"двигатель", "мощный", "жидкостный"},
			Metadata: map[string]*inventoryv1.Value{
				"тяга":    {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 3827000.0}},
				"топливо": {Value: &inventoryv1.Value_StringValue{StringValue: "керосин+кислород"}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"550e8400-e29b-41d4-a716-446655440002": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440002",
			Name:          "Крыло Delta-V",
			Description:   "Аэродинамическое крыло для атмосферного полета",
			Price:         2500000.0,
			StockQuantity: 8,
			Category:      inventoryv1.Category_CATEGORY_WING,
			Dimensions: &inventoryv1.Dimensions{
				Length: 1200.0,
				Width:  600.0,
				Height: 50.0,
				Weight: 850.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"крыло", "аэродинамика", "композит"},
			Metadata: map[string]*inventoryv1.Value{
				"материал":       {Value: &inventoryv1.Value_StringValue{StringValue: "углеродное волокно"}},
				"термостойкость": {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 1500.0}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"550e8400-e29b-41d4-a716-446655440003": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440003",
			Name:          "Криогенное топливо LH2",
			Description:   "Жидкий водород для ракетных двигателей",
			Price:         50000.0,
			StockQuantity: 150,
			Category:      inventoryv1.Category_CATEGORY_FUEL,
			Dimensions: &inventoryv1.Dimensions{
				Length: 100.0,
				Width:  100.0,
				Height: 200.0,
				Weight: 70.8,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Air Liquide",
				Country: "France",
				Website: "www.airliquide.com",
			},
			Tags: []string{"топливо", "криогенное", "водород"},
			Metadata: map[string]*inventoryv1.Value{
				"температура": {Value: &inventoryv1.Value_DoubleValue{DoubleValue: -253.0}},
				"чистота":     {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 99.9}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"550e8400-e29b-41d4-a716-446655440004": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440004",
			Name:          "Иллюминатор Космос-360",
			Description:   "Прочный иллюминатор для наблюдения в космосе",
			Price:         750000.0,
			StockQuantity: 12,
			Category:      inventoryv1.Category_CATEGORY_PORTHOLE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 60.0,
				Width:  60.0,
				Height: 15.0,
				Weight: 25.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Roscosmos",
				Country: "Russia",
				Website: "www.roscosmos.ru",
			},
			Tags: []string{"иллюминатор", "обзор", "прочный"},
			Metadata: map[string]*inventoryv1.Value{
				"материал_стекла": {Value: &inventoryv1.Value_StringValue{StringValue: "сапфировое стекло"}},
				"давление":        {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 101325.0}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"550e8400-e29b-41d4-a716-446655440005": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440005",
			Name:          "Двигатель Merlin 1D",
			Description:   "Компактный двигатель для первой ступени",
			Price:         1200000.0,
			StockQuantity: 25,
			Category:      inventoryv1.Category_CATEGORY_ENGINE,
			Dimensions: &inventoryv1.Dimensions{
				Length: 300.0,
				Width:  100.0,
				Height: 100.0,
				Weight: 630.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "SpaceX",
				Country: "USA",
				Website: "www.spacex.com",
			},
			Tags: []string{"двигатель", "компактный", "многоразовый"},
			Metadata: map[string]*inventoryv1.Value{
				"тяга":         {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 845000.0}},
				"многоразовый": {Value: &inventoryv1.Value_BoolValue{BoolValue: true}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
		"550e8400-e29b-41d4-a716-446655440006": {
			Uuid:          "550e8400-e29b-41d4-a716-446655440006",
			Name:          "Крыло Falcon Heavy",
			Description:   "Большое крыло для тяжелых ракет",
			Price:         4200000.0,
			StockQuantity: 4,
			Category:      inventoryv1.Category_CATEGORY_WING,
			Dimensions: &inventoryv1.Dimensions{
				Length: 1800.0,
				Width:  900.0,
				Height: 80.0,
				Weight: 1500.0,
			},
			Manufacturer: &inventoryv1.Manufacturer{
				Name:    "Blue Origin",
				Country: "USA",
				Website: "www.blueorigin.com",
			},
			Tags: []string{"крыло", "тяжелое", "стабилизация"},
			Metadata: map[string]*inventoryv1.Value{
				"грузоподъемность": {Value: &inventoryv1.Value_DoubleValue{DoubleValue: 63800.0}},
				"сертификат":       {Value: &inventoryv1.Value_StringValue{StringValue: "NASA-2024"}},
			},
			CreatedAt: now,
			UpdatedAt: now,
		},
	}
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Printf("Ошибка слушания tcp соединения на порту %d: %v\n", grpcPort, err)
		return
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("Ошибка закрытия tcp соединения на порту %d: %v\n", grpcPort, err)
		}
	}()

	s := grpc.NewServer()

	service := &inventoryService{
		storage: createTestData(),
		mu:      &sync.RWMutex{},
	}

	inventoryv1.RegisterInventoryServiceServer(s, service)

	reflection.Register(s)

	go func() {
		log.Printf("🚀 gRPC server listening on %d\n", grpcPort)
		err = s.Serve(lis)
		if err != nil {
			log.Printf("failed to serve: %v\n", err)
			return
		}
	}()
	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down gRPC server...")
	s.GracefulStop()
	log.Println("✅ Server stopped")
}
