package config

import (
	"fmt"
	"os"
)

type Config struct {
	DBURL     string
	JWTSecret string
	Port      string
}

func LoadConfig() (*Config, error) {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DB_URL environment variable is not set")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is not set")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DBURL:     dbURL,
		JWTSecret: jwtSecret,
		Port:      port,
	}, nil
}
