// iam/internal/service/dto/user_dto.go

package dto

// RegisterUserRequest DTO для регистрации пользователя
type RegisterUserRequest struct {
	Login               string
	Password            string // Открытый пароль
	Email               string
	NotificationMethods []NotificationMethodDTO
}

// RegisterUserResponse DTO ответа на регистрацию
type RegisterUserResponse struct {
	UserUUID string
}

// GetUserRequest DTO для получения пользователя
type GetUserRequest struct {
	UserUUID string
}

// GetUserResponse DTO пользователя
type GetUserResponse struct {
	UserUUID            string
	Login               string
	Email               string
	NotificationMethods []NotificationMethodDTO
	CreatedAt           string // ISO 8601
	UpdatedAt           string // ISO 8601
}

// UpdateUserRequest DTO для обновления пользователя
type UpdateUserRequest struct {
	UserUUID string

	// Опциональные поля (nil = не обновляем)
	Email               *string
	Password            *string // Открытый пароль
	NotificationMethods []NotificationMethodDTO
}

// NotificationMethodDTO DTO метода уведомления
type NotificationMethodDTO struct {
	ProviderName string // "telegram", "email", "sms"
	Target       string // chat_id, email, phone
}
