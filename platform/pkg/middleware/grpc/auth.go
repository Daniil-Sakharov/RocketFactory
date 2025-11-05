package grpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	authv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
	commonv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/common/v1"
)

const (
	SessionUUIDMetadataKey = "session-uuid"
)

type contextKey string

const (
	userContextKey        contextKey = "user"
	sessionUUIDContextKey contextKey = "session-uuid"
)

type AuthClient = authv1.AuthServiceClient

type AuthInterceptor struct {
	authClient AuthClient
}

func NewAuthInterceptor(authClient AuthClient) *AuthInterceptor {
	return &AuthInterceptor{authClient: authClient}
}

func (a *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		authCtx, err := a.authenticate(ctx)
		if err != nil {
			return nil, err
		}

		return handler(authCtx, req)
	}
}

func (a *AuthInterceptor) authenticate(ctx context.Context) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Error(codes.Unauthenticated, "missing metadata")
	}

	sessionUUIDs := md.Get(SessionUUIDMetadataKey)
	if len(sessionUUIDs) == 0 {
		return nil, status.Error(codes.Unavailable, "missing session-uuid in metadata")
	}

	sessionUUID := sessionUUIDs[0]
	if sessionUUID == "" {
		return nil, status.Error(codes.Unauthenticated, "empty session-uuid")
	}

	whoamiRes, err := a.authClient.Whoami(ctx, &authv1.WhoamiRequest{
		SessionUuid: sessionUUID,
	})
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, fmt.Sprintf("invalid session: %v", err))
	}

	authCtx := context.WithValue(ctx, userContextKey, whoamiRes.User)
	authCtx = context.WithValue(authCtx, sessionUUIDContextKey, sessionUUID)
	return authCtx, nil
}

func GetUserFromContext(ctx context.Context) (*commonv1.User, bool) {
	user, ok := ctx.Value(userContextKey).(*commonv1.User)
	return user, ok
}

func GetUserContextKey() contextKey {
	return userContextKey
}

func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	sessionUUID, ok := ctx.Value(sessionUUIDContextKey).(string)
	return sessionUUID, ok
}

func AddSessionUUIDToContext(ctx context.Context, sessionUUID string) context.Context {
	return context.WithValue(ctx, sessionUUIDContextKey, sessionUUID)
}

func ForwardSessionUUIDToGRPC(ctx context.Context) context.Context {
	sessionUUID, ok := GetSessionUUIDFromContext(ctx)
	if !ok || sessionUUID == "" {
		return ctx
	}

	return metadata.AppendToOutgoingContext(ctx, SessionUUIDMetadataKey, sessionUUID)
}
