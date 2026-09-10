package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_URL")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
