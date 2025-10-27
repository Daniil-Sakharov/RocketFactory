FROM golang:1.24.5-alpine AS builder

RUN apk add --no-cache git

WORKDIR /app

COPY go.work go.work.sum ./

COPY assembly/go.mod assembly/go.sum ./assembly/
COPY platform/go.mod platform/go.sum ./platform/
COPY shared/go.mod shared/go.sum ./shared/
COPY order/go.mod order/go.sum ./order/
COPY inventory/go.mod inventory/go.sum ./inventory/
COPY payment/go.mod payment/go.sum ./payment/
COPY notification/go.mod notification/go.sum ./notification/

# Скачиваем зависимости ПЕРЕД копированием исходников
# Этот слой будет закеширован, пока go.mod не изменится!
RUN go mod download

COPY assembly ./assembly
COPY platform ./platform
COPY shared ./shared
COPY order ./order
COPY inventory ./inventory
COPY payment ./payment
COPY notification ./notification

ADD https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/v0.4.37/grpc_health_probe-linux-amd64 ./grpc-health-probe

RUN chmod +x grpc-health-probe

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o app-inventory ./inventory/cmd/main.go

FROM alpine:3.21.3

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

USER appuser

COPY --from=builder /app/app-inventory ./app-inventory
COPY --from=builder /app/grpc-health-probe /bin/grpc-health-probe

EXPOSE 50051

ENTRYPOINT ["./app-inventory"]