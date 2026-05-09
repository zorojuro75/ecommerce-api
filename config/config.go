package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

type Config struct {
    ServerPort  string
    DatabaseURL string
    JWTSecret   string
}

func Load() *Config {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found — using system environment variables")
    }

    return &Config{
        ServerPort:  getEnv("PORT",         "8080"),
        DatabaseURL: getEnv("DATABASE_URL", ""),
        JWTSecret:   getEnv("JWT_SECRET",   "change-me-in-production"),
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}