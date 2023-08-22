package users

import (
	"time"
)

type User struct {
	Uuid                string    `json:"uuid" gorm:"primaryKey"`
	Username            string    `json:"username" gorm:"unique"`
	Email               string    `json:"email" gorm:"unique"`
	Password            string    `json:"password"`
	FirebaseDeviceToken string    `json:"firebase_device_token"`
	CreatedAt           time.Time `gorm:"type:timestamp;default:CURRENT_TIMESTAMP"`
	UpdatedAt           time.Time
}
