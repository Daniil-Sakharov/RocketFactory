package main

import (
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	apiv1 "github.com/Daniil-Sakharov/RocketFactory/order/internal/api/order/v1"
	inventoryClient "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc/inventory/v1"
	paymentClient "github.com/Daniil-Sakharov/RocketFactory/order/internal/client/grpc/payment/v1"
	"github.com/Daniil-Sakharov/RocketFactory/order/internal/migrator"
	orderRepo "github.com/Daniil-Sakharov/RocketFactory/order/internal/repository/order"
	orderService "github.com/Daniil-Sakharov/RocketFactory/order/internal/service/order"
	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
)

const (
	orderServicePort     = ":8080"
	orderServiceAddr     = "localhost:8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
)

// setupGRPCConnections создает подключения к внешним gRPC сервисам
func setupGRPCConnections() (*grpc.ClientConn, *grpc.ClientConn, error) {
	// Подключаемся к Inventory Service
	inventoryConn, err := grpc.NewClient(
		inventoryServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, nil, err
	}

	// Подключаемся к Payment Service
	paymentConn, err := grpc.NewClient(
		paymentServiceAddr,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		if closeErr := inventoryConn.Close(); closeErr != nil {
			log.Printf("Error closing inventory connection: %v", closeErr)
		}
		return nil, nil, err
	}

	return inventoryConn, paymentConn, nil
}

// run инициализирует и запускает Order Service
func run() error {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("failed to load .env file: %v\n", err)
		return err
	}
	dbURI := os.Getenv("DB_URI")
	migratorDir := os.Getenv("MIGRATIONS_DIR")
	db, err := sqlx.Connect("pgx", dbURI)
	if err != nil {
		log.Fatalf("❌ Failed to connect to PostgreSQL via sqlx: %v", err)
	}
	defer func() {
		if err = db.Close(); err != nil {
			log.Printf("failed to close db: %v", err)
		}
	}()
	if err = db.Ping(); err != nil {
		log.Printf("❌ Failed to ping PostgreSQL: %v", err)
	}

	mgrt := migrator.NewMigrator(db.DB, migratorDir)
	err = mgrt.Up()
	if err != nil {
		log.Printf("❌ Failed to make migrations: %v", err)
	}

	log.Println("🔌 Connecting to external services...")
	inventoryConn, paymentConn, err := setupGRPCConnections()
	if err != nil {
		return err
	}
	defer func() {
		if err := inventoryConn.Close(); err != nil {
			log.Printf("❌ Error closing inventory connection: %v", err)
		}
	}()
	defer func() {
		if err := paymentConn.Close(); err != nil {
			log.Printf("❌ Error closing payment connection: %v", err)
		}
	}()

	log.Printf("✅ Connected to Inventory Service at %s", inventoryServiceAddr)
	log.Printf("✅ Connected to Payment Service at %s", paymentServiceAddr)
	inventoryGRPCStub := inventoryV1.NewInventoryServiceClient(inventoryConn)
	paymentGRPCStub := paymentV1.NewPaymentServiceClient(paymentConn)

	inventoryGRPCClient := inventoryClient.NewClient(inventoryGRPCStub)
	paymentGRPCClient := paymentClient.NewClient(paymentGRPCStub)

	log.Println("✅ Client layer initialized")

	repository := orderRepo.NewRepository(db)

	log.Println("✅ Repository layer initialized (in-memory)")

	service := orderService.NewService(
		repository,
		inventoryGRPCClient,
		paymentGRPCClient,
	)

	log.Println("✅ Service layer initialized")

	apiHandler := apiv1.NewAPI(service)

	log.Println("✅ API layer initialized")

	server, err := orderV1.NewServer(apiHandler)
	if err != nil {
		return err
	}

	log.Println("✅ OpenAPI server created")

	httpMux := http.NewServeMux()

	// API endpoints
	httpMux.Handle("/api/", server)

	httpServer := &http.Server{
		Addr:         orderServicePort,
		Handler:      httpMux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return httpServer.ListenAndServe()
}

func main() {
	log.Println("🚀 Starting Order Service...")

	if err := run(); err != nil {
		log.Fatalf("❌ Application failed: %v", err)
	}
}
