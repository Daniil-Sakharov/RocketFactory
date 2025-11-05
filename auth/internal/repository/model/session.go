package model

import "time"

type Session struct {
	SessionUUID string
	UserUUID    string
	CreatedAt   time.Time
	ExpiresAt   time.Time
}

type SessionRedisView struct {
	SessionUUID string `redis:"session_uuid"`
	UserUUID    string `redis:"user_uuid"`
	CreatedAtNs int64  `redis:"created_at"`
	ExpiresAtNs int64  `redis:"expires_at"`
}
