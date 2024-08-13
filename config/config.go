package config

import (
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"log"
	"os"
)

var (
	PORT              = envDefault("BIND", ":8081")
	CONNECTION_STRING = envDefault("CONNECTION_STRING", "root:my-secret-pw@tcp(127.0.0.1:3306)/demo?charset=utf8mb4&parseTime=true")
	USERNAME          = envDefault("USERNAME", "admin")
	PASSWORD          = envDefault("PASSWORD", "admin")
)

func envDefault(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Failed to load .env file, using defaults:", err)
	}
}
