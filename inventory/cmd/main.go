package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	partAPIv1 "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/api/inventory/v1"
	partRepository "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/repository/part"
	partService "github.com/Daniil-Sakharov/RocketFactory/inventory/internal/service/part"
	inventoryv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
)

const grpcPort = 50051

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

	repo := partRepository.NewRepository()
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
