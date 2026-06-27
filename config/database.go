package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func LoadEnv() {
	_ = godotenv.Load()
}

func ConnectDB() (*gorm.DB, error) {
	LoadEnv()

	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		host := os.Getenv("DB_HOST")
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		name := os.Getenv("DB_NAME")
		port := os.Getenv("DB_PORT")
		if port == "" {
			port = "5432"
		}

		dsn = fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, name, port)
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	return db, nil
}
