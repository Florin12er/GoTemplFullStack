package models

import(
    "gorm.io/gorm"
)
// Notification model definition
type Notification struct {
    gorm.Model
    UserID     uint
    MessageID  uint
    Message    Message `gorm:"foreignKey:MessageID"`
    SenderID   uint
    Sender     User   `gorm:"foreignKey:SenderID"`
    IsRead     bool   `gorm:"default:false"`
}


