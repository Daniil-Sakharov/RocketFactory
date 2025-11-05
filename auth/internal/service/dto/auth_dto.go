// iam/internal/service/dto/auth_dto.go

package dto

// LoginRequest DTO для входа
type LoginRequest struct {
	Login    string
	Password string // Открытый пароль
}

// LoginResponse DTO ответа на вход
type LoginResponse struct {
	SessionUUID string
}

// WhoamiRequest DTO для проверки сессии
type WhoamiRequest struct {
	SessionUUID string
}

// WhoamiResponse DTO информации о текущем пользователе
type WhoamiResponse struct {
	UserUUID string
	Login    string
	Email    string
}
