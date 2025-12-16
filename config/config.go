package config

import (
	"os"
	"strconv"
)

type Config struct {
	LinkedInEmail        string
	LinkedInPassword     string
	DailyConnectionLimit int
	ConnectionNote       string
	MinDelay             int
	MaxDelay             int
}

func LoadConfig() *Config {
	dailyLimit, _ := strconv.Atoi(getEnv("DAILY_CONNECTION_LIMIT", "50"))
	minDelay, _ := strconv.Atoi(getEnv("MIN_DELAY_MS", "2000"))
	maxDelay, _ := strconv.Atoi(getEnv("MAX_DELAY_MS", "5000"))

	return &Config{
		LinkedInEmail:        getEnv("LINKEDIN_EMAIL", ""),
		LinkedInPassword:     getEnv("LINKEDIN_PASSWORD", ""),
		DailyConnectionLimit: dailyLimit,
		ConnectionNote:       getEnv("CONNECTION_NOTE", "Hi! I'd love to connect and exchange ideas."),
		MinDelay:             minDelay,
		MaxDelay:             maxDelay,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
