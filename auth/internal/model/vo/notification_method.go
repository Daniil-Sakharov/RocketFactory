package vo

import (
	"errors"
	"regexp"
)

var (
	ErrNotificationProviderEmpty   = errors.New("notification provider cannot be empty")
	ErrNotificationProviderInvalid = errors.New("invalid notification provider")
	ErrNotificationTargetEmpty     = errors.New("notification target cannot be empty")
	ErrNotificationTargetInvalid   = errors.New("invalid notification target format")
)

type NotificationMethod struct {
	providerName string
	target       string
}

func NewNotificationMethod(providerName, target string) (*NotificationMethod, error) {
	if providerName == "" {
		return nil, ErrNotificationProviderEmpty
	}

	allowedProvider := map[string]bool{
		"telegram": true,
		"email":    true,
		"sms":      true,
	}
	if !allowedProvider[providerName] {
		return nil, ErrNotificationProviderInvalid
	}

	if target == "" {
		return nil, ErrNotificationTargetEmpty
	}

	nm := &NotificationMethod{
		providerName: providerName,
		target:       target,
	}

	if err := nm.validateTarget(); err != nil {
		return nil, err
	}

	return nm, nil
}

func (nm *NotificationMethod) ProviderName() string {
	return nm.providerName
}

func (nm *NotificationMethod) Target() string {
	return nm.target
}

// Equals сравнивает два NotificationMethod
func (nm *NotificationMethod) Equals(other *NotificationMethod) bool {
	if other == nil {
		return false
	}
	return nm.providerName == other.providerName && nm.target == other.target
}

// validateTarget проверяет формат target в зависимости от провайдера
func (nm *NotificationMethod) validateTarget() error {
	switch nm.providerName {
	case "telegram":
		return nm.validateTelegramTarget()
	case "email":
		return nm.validateEmailTarget()
	case "sms":
		return nm.validateSMSTarget()
	}
	return nil
}

// validateTelegramTarget проверяет формат Telegram chat_id
func (nm *NotificationMethod) validateTelegramTarget() error {
	if len(nm.target) == 0 {
		return ErrNotificationTargetInvalid
	}
	return nil
}

// validateEmailTarget проверяет формат email
func (nm *NotificationMethod) validateEmailTarget() error {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !emailRegex.MatchString(nm.target) {
		return ErrNotificationTargetInvalid
	}
	return nil
}

// validateSMSTarget проверяет формат телефона
func (nm *NotificationMethod) validateSMSTarget() error {
	phoneRegex := regexp.MustCompile(`^\+[0-9]{10,15}$`)
	if !phoneRegex.MatchString(nm.target) {
		return ErrNotificationTargetInvalid
	}
	return nil
}
