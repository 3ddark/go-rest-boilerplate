package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost        string
	DBPort        string
	DBUser        string
	DBPassword    string
	DBName        string
	JWTSecret     string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		// .env dosyası yoksa hata verme, ortam değişkenlerinden okumaya devam et
		fmt.Println("No .env file found, reading from environment variables")
	}

	redisDB := 0 // Default DB
	redisDBStr := os.Getenv("REDIS_DB")
	if redisDBStr != "" {
		_, err := fmt.Sscanf(redisDBStr, "%d", &redisDB)
		if err != nil {
			return nil, fmt.Errorf("could not parse REDIS_DB: %w", err)
		}
	}

	return &Config{
		DBHost:        os.Getenv("DB_HOST"),
		DBPort:        os.Getenv("DB_PORT"),
		DBUser:        os.Getenv("DB_USER"),
		DBPassword:    os.Getenv("DB_PASSWORD"),
		DBName:        os.Getenv("DB_NAME"),
		JWTSecret:     os.Getenv("JWT_SECRET"),
		RedisAddr:     os.Getenv("REDIS_ADDR"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		RedisDB:       redisDB,
	}, nil
}
