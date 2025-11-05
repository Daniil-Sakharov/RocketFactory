package repository

import "fmt"

const (
	sessionKeyPrefix      = "auth:session:"
	userSessionsKeyPrefix = "auth:user_sessions:"
)

func SessionKey(sessionUUID string) string {
	return fmt.Sprintf("%s%s", sessionKeyPrefix, sessionUUID)
}

func UserSessionKey(userUUID string) string {
	return fmt.Sprintf("%s%s", sessionKeyPrefix, userUUID)
}
