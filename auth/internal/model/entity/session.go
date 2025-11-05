package entity

import "time"

type Session struct {
	sessionUUID string
	userUUID    string
	createdAt   time.Time
	expiresAt   time.Time
}

func NewSession(sessionUUID, userUUID string, ttl time.Duration) *Session {
	now := time.Now()
	return &Session{
		sessionUUID: sessionUUID,
		userUUID:    userUUID,
		createdAt:   now,
		expiresAt:   now.Add(ttl),
	}
}

func RestoreSession(sessionUUID, userUUID string, createdAt, expiredAt time.Time) *Session {
	return &Session{
		sessionUUID: sessionUUID,
		userUUID:    userUUID,
		createdAt:   createdAt,
		expiresAt:   expiredAt,
	}
}

func (s *Session) SessionUUID() string {
	return s.sessionUUID
}

func (s *Session) UserUUID() string {
	return s.userUUID
}

func (s *Session) CreatedAt() time.Time {
	return s.createdAt
}

func (s *Session) ExpiresAt() time.Time {
	return s.expiresAt
}

func (s *Session) IsExpired() bool {
	return time.Now().After(s.expiresAt)
}

func (s *Session) TTL() time.Duration {
	if s.IsExpired() {
		return 0
	}
	return time.Until(s.expiresAt)
}

func (s *Session) ShouldRefresh(threshold time.Duration) bool {
	return s.TTL() < threshold && s.TTL() > 0
}
