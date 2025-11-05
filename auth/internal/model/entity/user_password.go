package entity

import (
	"time"

	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model"
	"github.com/Daniil-Sakharov/RocketFactory/auth/internal/model/vo"
)

func (u *User) UpdatePassword(newPassword *vo.Password) error {
	if u.password.Hash() == newPassword.Hash() {
		return model.ErrPasswordSameAsCurrent
	}

	u.password = newPassword
	u.updatedAt = time.Now()

	return nil
}
