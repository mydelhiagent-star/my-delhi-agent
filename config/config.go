package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	MongoURI  string
	MongoDB   string
	JWTSecret string
	Port      string
	RedisURI  string
	RedisUsername string
	RedisPassword string
	AdminEmail string
	AdminPassword string
	CloudflareAccountID       string
	CloudflareAccessKeyID     string
	CloudflareAccessKeySecret string
	CloudflareBucketName      string
}

func LoadConfig() Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	return Config{
		MongoURI:  os.Getenv("MONGO_URI"),
		MongoDB:   os.Getenv("MONGO_DB"),
		JWTSecret: os.Getenv("JWT_SECRET"),
		Port:      os.Getenv("PORT"),
		RedisURI:  os.Getenv("REDIS_URI"),
		RedisUsername: os.Getenv("REDIS_USERNAME"),
		RedisPassword: os.Getenv("REDIS_PASSWORD"),
		AdminEmail: os.Getenv("ADMIN_EMAIL"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
		CloudflareAccountID: os.Getenv("CLOUDFARE_ACCOUNT_ID"),
		CloudflareAccessKeyID: os.Getenv("CLOUDFARE_ACCESS_KEY_ID"),
		CloudflareAccessKeySecret: os.Getenv("CLOUDFARE_ACCESS_KEY_SECRET"),
		CloudflareBucketName: os.Getenv("CLOUDFARE_BUCKET_NAME"),
	}
}
