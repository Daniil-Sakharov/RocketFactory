package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"
)

type User struct {
	UserUUID            string                `db:"user_uuid"`
	Login               string                `db:"login"`
	PasswordHash        string                `db:"password_hash"`
	Email               string                `db:"email"`
	NotificationMethods NotificationMethodsDB `db:"notification_methods"`
	CreatedAt           time.Time             `db:"created_at"`
	UpdatedAt           time.Time             `db:"updated_at"`
}

type NotificationMethodsDB []NotificationMethodDB

type NotificationMethodDB struct {
	ProviderName string `json:"provider_name"`
	Target       string `json:"target"`
}

func (nm *NotificationMethodsDB) Value() (driver.Value, error) {
	if nm == nil {
		return json.Marshal([]NotificationMethodDB{})
	}
	return json.Marshal(nm)
}

func (nm *NotificationMethodsDB) Scan(value interface{}) error {
	if value == nil {
		*nm = []NotificationMethodDB{}
		return nil
	}

	bytes, ok := value.([]byte)
	if !ok {
		return nil
	}

	return json.Unmarshal(bytes, nm)
}
