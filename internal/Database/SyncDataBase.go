package database

import (
	"log"
	"GoMessageApp/internal/models"
)

// SyncDatabase migrates the database and handles errors
func SyncDatabase() {
	if err := DB.AutoMigrate(&models.User{}, &models.Message{}, &models.Notification{}); err != nil {
		log.Fatalf("failed to migrate database: %v", err)
	}
}

