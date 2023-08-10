package internal

import "github.com/joho/godotenv"

// LoadEnv loads env vars from .env file
func LoadEnv() error {
	return godotenv.Load()
}
