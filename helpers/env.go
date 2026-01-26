package helpers

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

// Fungsi untuk load .env dan ambil datanya
func GetConfig(key string) string {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	return os.Getenv(key)
}

// var DB_DRIVER = "mysql"
// var DB_DRIVER = "postgre"

// var DB_HOST = "127.0.0.1"
// var DB_PORT = "5432"
// // var DB_PORT = "3306"
// var DB_NAME = "go_gin_gorm"
// var DB_USER = "postgres"
// var DB_PASSWORD = "teazet"
