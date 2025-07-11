package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() error {
	return godotenv.Load("config.env")
}

// GetDBConfig returns database configuration
func GetDBConfig() map[string]string {
	return map[string]string{
		"host":     os.Getenv("DB_HOST"),
		"port":     os.Getenv("DB_PORT"),
		"user":     os.Getenv("DB_USER"),
		"password": os.Getenv("DB_PASSWORD"),
		"dbname":   os.Getenv("DB_NAME"),
	}
}

// GetJWTConfig returns JWT configuration
func GetJWTConfig() map[string]string {
	return map[string]string{
		"secret": os.Getenv("JWT_SECRET"),
		"expiry": os.Getenv("JWT_EXPIRY"),
	}
}

// GetUploadConfig returns file upload configuration
func GetUploadConfig() map[string]interface{} {
	maxFileSize, _ := strconv.ParseInt(os.Getenv("MAX_FILE_SIZE"), 10, 64)
	return map[string]interface{}{
		"upload_path":   os.Getenv("UPLOAD_PATH"),
		"max_file_size": maxFileSize,
	}
} 