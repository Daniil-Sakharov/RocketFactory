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

	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"

	apiv1 "github.com/Daniil-Sakharov/RocketFactory/payment/internal/api/payment/v1"
	paymentService "github.com/Daniil-Sakharov/RocketFactory/payment/internal/service/payment"
)

const grpcPort = 50052

func main() {
	// Создаем TCP listener
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatalf("Ошибка слушания tcp соединения на порту %d: %v\n", grpcPort, err)
	}
	defer func() {
		if err := lis.Close(); err != nil {
			log.Printf("Ошибка закрытия tcp соединения на порту %d: %v\n", grpcPort, err)
		}
	}()

	// Создаем Service слой
	service := paymentService.New()

	// Создаем API слой
	api := apiv1.New(service)

	// Создаем gRPC сервер
	grpcServer := grpc.NewServer()

	// Регистрируем API
	paymentv1.RegisterPaymentServiceServer(grpcServer, api)

	// Включаем рефлексию для grpcurl
	reflection.Register(grpcServer)

	// Запускаем сервер в горутине
	go func() {
		log.Printf("🚀 PaymentService gRPC server listening on port %d\n", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v\n", err)
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down PaymentService gRPC server...")
	grpcServer.GracefulStop()
	log.Println("✅ PaymentService stopped")
}
