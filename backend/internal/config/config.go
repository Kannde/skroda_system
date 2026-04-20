package config

import (
	"os"
	"strconv"
)

type Config struct {
	Port                  string
	GinMode               string
	DatabaseURL           string
	RedisURL              string
	JWTSecret             string
	JWTExpiry             string
	MoMoAPIKey            string
	MoMoAPISecret         string
	MoMoCallbackURL       string
	ResendAPIKey          string
	SMSAPIKey             string
	SMSAPIURL             string
	AppURL                string
	InspectionPeriodHours int
}

func Load() *Config {
	inspectionHours, _ := strconv.Atoi(getEnv("INSPECTION_PERIOD_HOURS", "48"))
	return &Config{
		Port:                  getEnv("PORT", "8080"),
		GinMode:               getEnv("GIN_MODE", "debug"),
		DatabaseURL:           getEnv("DATABASE_URL", ""),
		RedisURL:              getEnv("REDIS_URL", "redis://localhost:6379"),
		JWTSecret:             getEnv("JWT_SECRET", ""),
		JWTExpiry:             getEnv("JWT_EXPIRY", "24h"),
		MoMoAPIKey:            getEnv("MOMO_API_KEY", ""),
		MoMoAPISecret:         getEnv("MOMO_API_SECRET", ""),
		MoMoCallbackURL:       getEnv("MOMO_CALLBACK_URL", ""),
		ResendAPIKey:          getEnv("RESEND_API_KEY", ""),
		SMSAPIKey:             getEnv("SMS_API_KEY", ""),
		SMSAPIURL:             getEnv("SMS_API_URL", ""),
		AppURL:                getEnv("APP_URL", "http://localhost:5173"),
		InspectionPeriodHours: inspectionHours,
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
