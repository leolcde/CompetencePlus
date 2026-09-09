package db

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect() (*gorm.DB, error) {
	dsn := os.Getenv("DB_URL") // provided by docker-compose (host = "postgres")
	if dsn == "" {
		dsn = os.Getenv("POSTGRES_URL") // local dev fallback (host = "localhost")
	}
	return gorm.Open(postgres.Open(dsn), &gorm.Config{})
}
