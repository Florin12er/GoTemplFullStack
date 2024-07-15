package models

import(
    "time"
    "gorm.io/gorm"
)
// Message model definition
type Message struct {
	ID         uint           `gorm:"primaryKey"`
	CreatedAt  time.Time      `gorm:"not null"`
	UpdatedAt  time.Time      `gorm:"not null"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	SenderID   uint           `gorm:"not null"` // Foreign key to reference User
	ReceiverID uint           `gorm:"not null"` // Foreign key to reference the receiving User
	Text       string         `gorm:"not null"`
	Sender     User           `gorm:"foreignKey:SenderID"`
	Receiver   User           `gorm:"foreignKey:ReceiverID"`
	Picture    string
}

