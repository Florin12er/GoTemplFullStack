package models

import (
	"time"

	"gorm.io/gorm"
)

// User model definition
type User struct {
	ID             uint           `gorm:"primaryKey"`
	CreatedAt      time.Time      `gorm:"not null"`
	UpdatedAt      time.Time      `gorm:"not null"`
	DeletedAt      gorm.DeletedAt `gorm:"index"`
	UserName       string         `gorm:"unique;not null"`
	FullName       string         `gorm:"not null"`
	Email          string         `gorm:"unique;not null"`
	Password       string         `gorm:"not null"`
	ProfilePicture string
	Messages       []Message `gorm:"foreignKey:SenderID"`
}

