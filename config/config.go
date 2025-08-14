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
	}
}
