package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partAPIv1 "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/api/inventory/v1"
	partRepository "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/part"
	partService "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service/part"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

func main() {
	ctx := context.Background()

	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Ошибка загружения .env файла")
		return
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Printf("Ошибка: переменная окружения MONGO_URI не установлена")
		return
	}

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Printf("Ошибка подключения к MongoDB: %v\n", err)
		return
	}

	defer func() {
		if cerr := client.Disconnect(ctx); cerr != nil {
			log.Printf("Ошибка закрытия соединения MongoDB")
		}
	}()

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Printf("MongoDB недоступен, ошибка ping: %v\n", err)
		return
	}
	log.Printf("✅ Успешное подключение к MongoDB")

	MongoDB := client.Database("inventory-service")

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

	repo := partRepository.NewRepository(MongoDB)
	repo.InitTestData()
	service := partService.NewService(repo)
	api := partAPIv1.NewAPI(service)

	inventoryv1.RegisterInventoryServiceServer(s, api)

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
