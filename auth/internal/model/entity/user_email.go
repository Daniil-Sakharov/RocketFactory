package entity

import (
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
)

func (u *User) UpdateEmail(newEmail *vo.Email) error {
	if u.email.Value() == newEmail.Value() {
		return model.ErrEmailSameAsCurrent
	}

	u.email = newEmail
	u.updatedAt = time.Now()

	return nil
}
