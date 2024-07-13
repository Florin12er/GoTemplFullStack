package models

import (
	"time"

	"gorm.io/gorm"
)

// User model definition
type User struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time      `gorm:"not null"`
	UpdatedAt time.Time      `gorm:"not null"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	UserName  string         `gorm:"unique;not null"`
	FullName  string         `gorm:"not null"`
	Email     string         `gorm:"unique;not null"`
	Password  string         `gorm:"not null"`
}

// Message model definition
type Message struct {
	ID         uint           `gorm:"primaryKey"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	SenderID   uint           // Foreign key to reference User
	ReceiverID uint           // Foreign key to reference the receiving User
	Text       string         `gorm:"not null"`
	Sender     User           `gorm:"foreignKey:SenderID"`
	Receiver   User           `gorm:"foreignKey:ReceiverID"`
}

