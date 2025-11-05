package http

import (
	"context"
	"net/http"

	grpcAuth "github.com/Daniil-Sakharov/RocketFactory/platform/pkg/middleware/grpc"
	authv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/auth/v1"
	commonv1 "github.com/Daniil-Sakharov/RocketFactory/shared/pkg/proto/common/v1"
)

const SessionUUIDHeader = "X-Session-Uuid"

type AuthClient = authv1.AuthServiceClient

type AuthMiddleware struct {
	authClient AuthClient
}

func NewAuthMiddleware(authClient AuthClient) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: authClient,
	}
}

func (m *AuthMiddleware) Handle(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionUUID := r.Header.Get(SessionUUIDHeader)
		if sessionUUID == "" {
			writeErrorResponse(w, http.StatusUnauthorized, "MISSING_SESSION", "Authentication required")
			return
		}

		whoamiRes, err := m.authClient.Whoami(r.Context(), &authv1.WhoamiRequest{
			SessionUuid: sessionUUID,
		})
		if err != nil {
			writeErrorResponse(w, http.StatusUnauthorized, "INVALID_SESSION", "Authentication failed")
			return
		}

		ctx := r.Context()
		ctx = grpcAuth.AddSessionUUIDToContext(ctx, sessionUUID)
		ctx = context.WithValue(ctx, grpcAuth.GetUserContextKey(), whoamiRes.User)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*commonv1.User, bool) {
	return grpcAuth.GetUserFromContext(ctx)
}

func GetSessionUUIDFromContext(ctx context.Context) (string, bool) {
	return grpcAuth.GetSessionUUIDFromContext(ctx)
}
