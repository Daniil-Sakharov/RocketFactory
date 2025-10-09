package main

import (
	"context"
	"log"
	"net/http"
	"sync"

	orderV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/openapi/order/v1"
	inventoryV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/inventory/v1"
	paymentV1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/payment/v1"
	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	orderServicePort     = ":8080"
	orderServiceAddr     = "localhost:8080"
	inventoryServiceAddr = "localhost:50051"
	paymentServiceAddr   = "localhost:50052"
)

// Структура заказа для внутреннего хранения
type Order struct {
	UUID            uuid.UUID              `json:"order_uuid"`
	UserUUID        uuid.UUID              `json:"user_uuid"`
	PartUUIDs       []uuid.UUID            `json:"part_uuids"`
	TotalPrice      float64                `json:"total_price"`
	Status          orderV1.OrderStatus    `json:"status"`
	TransactionUUID *uuid.UUID             `json:"transaction_uuid,omitempty"`
	PaymentMethod   *orderV1.PaymentMethod `json:"payment_method,omitempty"`
}

// Хранилище заказов
type OrderStorage struct {
	mu     *sync.RWMutex
	orders map[uuid.UUID]*Order
}

// Основной сервис
type OrderService struct {
	storage         *OrderStorage
	inventoryClient inventoryV1.InventoryServiceClient
	paymentClient   paymentV1.PaymentServiceClient
}

func NewOrderStorage() *OrderStorage {
	return &OrderStorage{
		mu:     &sync.RWMutex{},
		orders: make(map[uuid.UUID]*Order),
	}
}

func (s *OrderService) GetOrder(ctx context.Context, params orderV1.GetOrderParams) (orderV1.GetOrderRes, error) {
	s.storage.mu.RLock()
	defer s.storage.mu.RUnlock()

	// Ищем заказ в storage
	order, exists := s.storage.orders[params.OrderUUID]
	if !exists {
		// Возвращаем 404 Not Found
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: "Order not found",
		}, nil
	}

	// Конвертируем внутреннюю структуру в OpenAPI response
	response := &orderV1.GetOrderResponse{
		OrderUUID:  order.UUID,
		UserUUID:   order.UserUUID,
		PartUuids:  order.PartUUIDs,
		TotalPrice: order.TotalPrice,
		Status:     order.Status,
	}

	// Добавляем опциональные поля если заказ оплачен
	if order.TransactionUUID != nil {
		response.TransactionUUID.SetTo(*order.TransactionUUID)
	}
	if order.PaymentMethod != nil {
		response.PaymentMethod.SetTo(*order.PaymentMethod)
	}

	return response, nil
}

// Заглушки для остальных методов (ты их реализуешь сам)
func (s *OrderService) CreateOrder(ctx context.Context, req *orderV1.CreateOrderRequest) (orderV1.CreateOrderRes, error) {
	request := &inventoryV1.ListPartsRequest{
		Filter: &inventoryV1.PartsFilter{
			Uuids: make([]string, len(req.PartUuids)),
		},
	}

	for i, partUUID := range req.PartUuids {
		request.Filter.Uuids[i] = partUUID.String()
	}

	response, err := s.inventoryClient.ListParts(ctx, request)
	if err != nil {
		return &orderV1.BadGatewayError{
			Error:   "INVENTORY_SERVICE_ERROR",
			Message: "Failed to check parts availability",
		}, nil
	}

	if len(response.Parts) != len(req.PartUuids) {
		return &orderV1.ValidationError{
			Error:   "PARTS_NOT_FOUND",
			Message: "Some requested parts are not available",
		}, nil
	}

	foundUUIDs := make(map[string]bool)
	for _, part := range response.Parts {
		foundUUIDs[part.Uuid] = true
	}

	for _, requestedUUID := range req.PartUuids {
		if !foundUUIDs[requestedUUID.String()] {
			return &orderV1.ValidationError{
				Error:   "PARTS_NOT_FOUND",
				Message: "Some requested parts are not available",
			}, nil
		}
	}

	var totalPrice float64
	for _, part := range response.Parts {
		totalPrice += part.Price
	}

	newOrderUUID := uuid.New()
	s.storage.mu.Lock()
	s.storage.orders[newOrderUUID] = &Order{
		UUID:            newOrderUUID,
		UserUUID:        req.UserUUID,
		PartUUIDs:       req.PartUuids,
		TotalPrice:      totalPrice,
		Status:          orderV1.OrderStatusPENDINGPAYMENT,
		TransactionUUID: nil,
		PaymentMethod:   nil,
	}
	s.storage.mu.Unlock()

	resp := &orderV1.CreateOrderResponse{
		OrderUUID:  newOrderUUID,
		TotalPrice: totalPrice,
	}

	return resp, nil
}

func (s *OrderService) PayOrder(ctx context.Context, req *orderV1.PayOrderRequest, params orderV1.PayOrderParams) (orderV1.PayOrderRes, error) {
	s.storage.mu.RLock()
	order, exist := s.storage.orders[params.OrderUUID]
	s.storage.mu.RUnlock()
	if !exist {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: "Order not found",
		}, nil
	}
	if order.Status != orderV1.OrderStatusPENDINGPAYMENT {
		return &orderV1.ConflictError{
			Error:   "INVALID_STATUS",
			Message: "Order is not pending payment",
		}, nil
	}

	paymentMethod := convertPaymentMethod(req.PaymentMethod)
	request := paymentV1.PayOrderRequest{
		OrderUuid:     order.UUID.String(),
		UserUuid:      order.UserUUID.String(),
		PaymentMethod: paymentMethod,
	}
	response, err := s.paymentClient.PayOrder(ctx, &request)
	if err != nil {
		return &orderV1.BadGatewayError{
			Error:   "PAYMENT_SERVICE_ERROR",
			Message: "Error to process payment",
		}, nil
	}

	transactionUUID, err := uuid.Parse(response.TransactionUuid)
	if err != nil {
		return &orderV1.ValidationError{
			Error:   "PARSE_ERROR",
			Message: "Error to parse string to UUID",
		}, nil
	}

	resp := orderV1.PayOrderResponse{
		TransactionUUID: transactionUUID,
	}

	s.storage.mu.Lock()
	order.Status = orderV1.OrderStatusPAID
	order.TransactionUUID = &transactionUUID
	order.PaymentMethod = &req.PaymentMethod
	s.storage.mu.Unlock()

	return &resp, nil
}

func convertPaymentMethod(apiMethod orderV1.PaymentMethod) paymentV1.PaymentMethod {
	switch apiMethod {
	case orderV1.PaymentMethodCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CARD
	case orderV1.PaymentMethodSBP:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_SBP
	case orderV1.PaymentMethodCREDITCARD:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_CREDIT_CARD
	case orderV1.PaymentMethodINVESTORMONEY:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_INVESTOR_MONEY
	default:
		return paymentV1.PaymentMethod_PAYMENT_METHOD_UNSPECIFIED
	}
}

func (s *OrderService) CancelOrder(ctx context.Context, params orderV1.CancelOrderParams) (orderV1.CancelOrderRes, error) {
	order, exist := s.storage.orders[params.OrderUUID]
	if !exist {
		return &orderV1.NotFoundError{
			Error:   "NOT_FOUND",
			Message: "Order not found",
		}, nil
	}
	if order.Status == orderV1.OrderStatusPAID {
		return &orderV1.ConflictError{
			Error:   "CONFLICT",
			Message: "Order already paid",
		}, nil
	}
	s.storage.mu.Lock()
	defer s.storage.mu.Unlock()
	order.Status = orderV1.OrderStatusCANCELLED
	return &orderV1.CancelOrderNoContent{}, nil
}

func main() {
	// Подключение к gRPC сервисам
	inventoryConn, err := grpc.NewClient(inventoryServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to inventory service: %v", err)
	}
	defer inventoryConn.Close()

	paymentConn, err := grpc.NewClient(paymentServiceAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to payment service: %v", err)
	}
	defer paymentConn.Close()

	service := &OrderService{
		storage:         NewOrderStorage(),
		inventoryClient: inventoryV1.NewInventoryServiceClient(inventoryConn),
		paymentClient:   paymentV1.NewPaymentServiceClient(paymentConn),
	}

	server, err := orderV1.NewServer(service)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	fileServer := http.FileServer(http.Dir("api"))
	log.Printf("📁 Serving static files from: api")

	httpMux := http.NewServeMux()

	httpMux.Handle("/api/", server)

	httpMux.Handle("/swagger-ui.html", fileServer)
	httpMux.Handle("/swagger.json", fileServer)

	httpMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			http.Redirect(w, r, "/swagger-ui.html", http.StatusMovedPermanently)
			return
		}
		fileServer.ServeHTTP(w, r)
	})

	log.Printf("🚀 OrderService HTTP server listening on %s", orderServiceAddr)
	log.Printf("📖 API Documentation available at: http://%s/swagger-ui.html", orderServiceAddr)
	log.Printf("🔗 Direct API access: http://%s/api/v1/orders", orderServiceAddr)

	if err := http.ListenAndServe(orderServicePort, httpMux); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
