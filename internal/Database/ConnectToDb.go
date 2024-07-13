package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectToDb() {
    var err error
    // Update the DSN to include sslmode and sslrootcert parameters
    dsn := os.Getenv("DB")
    DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})

    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
}

