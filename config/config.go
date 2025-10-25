package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string 
	DBString string 
}

// Loads the configuration from an .env variable 
func LoadConfig() (*Config, error) {
	config := Config{}
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get the env values for port and MySQL db host
	port := getVal("PORT", ":4000")
	dbString := getVal("DB_STRING", "")

	config.Port = port
	config.DBString = dbString

	return &config, err
}

func getVal(key, defaultValue string) string{
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}