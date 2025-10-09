package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	paymentv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

const grpcPort = 50052

type paymentService struct {
	paymentv1.UnimplementedPaymentServiceServer
}

func (s *paymentService) PayOrder(_ context.Context, _ *paymentv1.PayOrderRequest) (*paymentv1.PayOrderResponse, error) {
	transactionUUID := uuid.NewString()
	log.Printf("Оплата прошла успешно, transaction_uuid: %s", transactionUUID)
	return &paymentv1.PayOrderResponse{
		TransactionUuid: transactionUUID,
	}, nil
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

	service := &paymentService{}

	paymentv1.RegisterPaymentServiceServer(s, service)

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
