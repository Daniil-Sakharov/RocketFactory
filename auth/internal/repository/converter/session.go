package converter

import (
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/entity"
	repoModel "github.com/Daniil-Sakharov/RocketFactory/auth/internal/repository/model"
)

func EntityToRepositorySession(domainSession *entity.Session) *repoModel.Session {
	return &repoModel.Session{
		SessionUUID: domainSession.SessionUUID(),
		UserUUID:    domainSession.UserUUID(),
		CreatedAt:   domainSession.CreatedAt(),
		ExpiresAt:   domainSession.ExpiresAt(),
	}
}

func RepositoryToEntitySession(repoSession *repoModel.Session) *entity.Session {
	return entity.RestoreSession(
		repoSession.SessionUUID,
		repoSession.UserUUID,
		repoSession.CreatedAt,
		repoSession.ExpiresAt,
	)
}

func SessionToRedisView(session *repoModel.Session) *repoModel.SessionRedisView {
	return &repoModel.SessionRedisView{
		SessionUUID: session.SessionUUID,
		UserUUID:    session.UserUUID,
		CreatedAtNs: session.CreatedAt.UnixNano(),
		ExpiresAtNs: session.ExpiresAt.UnixNano(),
	}
}

func SessionFromRedisView(redisView *repoModel.SessionRedisView) *repoModel.Session {
	return &repoModel.Session{
		SessionUUID: redisView.SessionUUID,
		UserUUID:    redisView.UserUUID,
		CreatedAt:   time.Unix(0, redisView.CreatedAtNs),
		ExpiresAt:   time.Unix(0, redisView.ExpiresAtNs),
	}
}
