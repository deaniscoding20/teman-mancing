package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Config represents all application runtime configurations.
type Config struct {
	AppEnv         string
	AppPort        int
	DatabaseURL    string
	DBHost         string
	DBPort         int
	DBUser         string
	DBPassword     string
	DBName         string
	DBSSLMode      string
	JWTSecret      string
	JWTExpiration  time.Duration
	MinIOEndpoint  string
	MinIOAccessKey string
	MinIOSecretKey string
	MinIOUseSSL    bool
	MinIOBuckets   MinIOBuckets
	AllowedOrigins []string
	RedisURL       string
	SMTPHost       string
	SMTPPort       int
	SMTPUser       string
	SMTPPassword   string
	SMTPUseSSL     bool
	SMTPSender     string
	SMTPSenderName string
	AppBaseURL     string

	GoogleClientID     string
	GoogleClientSecret string
	GoogleRedirectURL  string
}

// MinIOBuckets specifies bucket names.
type MinIOBuckets struct {
	Avatars string
	Spots   string
	Events  string
}

// Load loads configuration from environment variables with sensible defaults.
func Load() *Config {
	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnvAsInt("APP_PORT", 8080),
		DatabaseURL:    getEnv("DATABASE_URL", "postgres://teman_mancing_user:teman_mancing_secret@localhost:5432/teman_mancing?sslmode=disable"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnvAsInt("DB_PORT", 5432),
		DBUser:         getEnv("DB_USER", "teman_mancing_user"),
		DBPassword:     getEnv("DB_PASSWORD", "teman_mancing_secret"),
		DBName:         getEnv("DB_NAME", "teman_mancing"),
		DBSSLMode:      getEnv("DB_SSLMODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "teman-mancing-secret-key-change-in-production-min-32-chars"),
		JWTExpiration:  time.Duration(getEnvAsInt("JWT_EXPIRATION_HOURS", 24)) * time.Hour,
		MinIOEndpoint:  getEnv("MINIO_ENDPOINT", "localhost:9000"),
		MinIOAccessKey: getEnv("MINIO_ACCESS_KEY", "minioadmin"),
		MinIOSecretKey: getEnv("MINIO_SECRET_KEY", "minioadminsecret"),
		MinIOUseSSL:    getEnvAsBool("MINIO_USE_SSL", false),
		MinIOBuckets: MinIOBuckets{
			Avatars: getEnv("MINIO_BUCKET_AVATARS", "avatars"),
			Spots:   getEnv("MINIO_BUCKET_SPOTS", "fishing-spots"),
			Events:  getEnv("MINIO_BUCKET_EVENTS", "event-covers"),
		},
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"http://localhost:3080", "http://127.0.0.1:3080", "http://localhost:3000", "http://127.0.0.1:3000"}),
		RedisURL:       getEnv("REDIS_URL", "redis://localhost:6379"),
		SMTPHost:       getEnv("SMTP_HOST", "smtp.sumopod.com"),
		SMTPPort:       getEnvAsInt("SMTP_PORT", 465),
		SMTPUser:       getEnv("SMTP_USER", "cmw0i63e310in08odzvjptyy7"),
		SMTPPassword:   getEnv("SMTP_PASSWORD", ""),
		SMTPUseSSL:     getEnvAsBool("SMTP_USE_SSL", true),
		SMTPSender:     getEnv("SMTP_SENDER", "noreply@sumopod.com"),
		SMTPSenderName: getEnv("SMTP_SENDER_NAME", "Teman Mancing"),
		AppBaseURL:     getEnv("APP_BASE_URL", "https://temanmancing.my.id"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
		GoogleRedirectURL:  getEnv("GOOGLE_REDIRECT_URL", "https://temanmancing.my.id/api/v1/auth/google/callback"),
	}
}

func (c *Config) IsDevelopment() bool {
	return strings.ToLower(c.AppEnv) == "development"
}

func (c *Config) IsProduction() bool {
	return strings.ToLower(c.AppEnv) == "production"
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valStr := getEnv(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsBool(key string, defaultVal bool) bool {
	valStr := getEnv(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}

func getEnvAsSlice(key string, defaultVal []string) []string {
	valStr := getEnv(key, "")
	if valStr == "" {
		return defaultVal
	}
	items := strings.Split(valStr, ",")
	res := make([]string, 0, len(items))
	for _, item := range items {
		trimmed := strings.TrimSpace(item)
		if trimmed != "" {
			res = append(res, trimmed)
		}
	}
	return res
}
