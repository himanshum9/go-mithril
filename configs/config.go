package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all application configuration
type Config struct {
	Database DatabaseConfig
	AWS      AWSConfig
	Kafka    KafkaConfig
	Streaming StreamingConfig
	Session   SessionConfig
	Security  SecurityConfig
	Logging   LoggingConfig
	Environment EnvironmentConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string
	SSLMode  string
}

type AWSConfig struct {
	Region           string
	CognitoUserPoolID string
	CognitoAppClientID string
	CognitoAppClientSecret string
}

type KafkaConfig struct {
	Broker   string
	Topic    string
	GroupID  string
}

type StreamingConfig struct {
	Endpoint      string
	RetryInterval int
	MaxRetries    int
	WebSocketPort int
}

type SessionConfig struct {
	DurationSeconds      int
	SubmissionIntervalSeconds int
}

type SecurityConfig struct {
	JWTSecret        string
	CORSAllowedOrigins string
}

type LoggingConfig struct {
	Level  string
	Format string
}

type EnvironmentConfig struct {
	Environment string
	Debug       bool
}

// Load loads configuration from environment variables
func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "your_db_user"),
			Password: getEnv("DB_PASSWORD", "your_db_password"),
			Name:     getEnv("DB_NAME", "multi_tenant_db"),
			SSLMode:  getEnv("DB_SSL_MODE", "disable"),
		},
		AWS: AWSConfig{
			Region:           getEnv("AWS_REGION", "us-west-2"),
			CognitoUserPoolID: getEnv("COGNITO_USER_POOL_ID", "your_user_pool_id"),
			CognitoAppClientID: getEnv("COGNITO_APP_CLIENT_ID", "your_app_client_id"),
			CognitoAppClientSecret: getEnv("COGNITO_APP_CLIENT_SECRET", "your_app_client_secret"),
		},
		Kafka: KafkaConfig{
			Broker:  getEnv("KAFKA_BROKER", "localhost:9092"),
			Topic:   getEnv("KAFKA_TOPIC", "location-stream"),
			GroupID: getEnv("KAFKA_GROUP_ID", "location-service-group"),
		},
		Streaming: StreamingConfig{
			Endpoint:      getEnv("STREAMING_ENDPOINT", "http://third-party-streaming-endpoint"),
			RetryInterval: getEnvAsInt("STREAMING_RETRY_INTERVAL", 5),
			MaxRetries:    getEnvAsInt("STREAMING_MAX_RETRIES", 3),
			WebSocketPort: getEnvAsInt("STREAMING_WEBSOCKET_PORT", 8084),
		},
		Session: SessionConfig{
			DurationSeconds:      getEnvAsInt("SESSION_DURATION_SECONDS", 600),
			SubmissionIntervalSeconds: getEnvAsInt("SUBMISSION_INTERVAL_SECONDS", 30),
		},
		Security: SecurityConfig{
			JWTSecret:        getEnv("JWT_SECRET", "your_jwt_secret_key_here"),
			CORSAllowedOrigins: getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:8080"),
		},
		Logging: LoggingConfig{
			Level:  getEnv("LOG_LEVEL", "info"),
			Format: getEnv("LOG_FORMAT", "json"),
		},
		Environment: EnvironmentConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
			Debug:       getEnvAsBool("DEBUG", true),
		},
	}
}

// GetDatabaseURL returns the database connection string
func (c *Config) GetDatabaseURL() string {
	return "postgres://" + c.Database.User + ":" + c.Database.Password + "@" + c.Database.Host + ":" + strconv.Itoa(c.Database.Port) + "/" + c.Database.Name + "?sslmode=" + c.Database.SSLMode
}

// GetSessionDuration returns session duration as time.Duration
func (c *Config) GetSessionDuration() time.Duration {
	return time.Duration(c.Session.DurationSeconds) * time.Second
}

// GetSubmissionInterval returns submission interval as time.Duration
func (c *Config) GetSubmissionInterval() time.Duration {
	return time.Duration(c.Session.SubmissionIntervalSeconds) * time.Second
}

// Helper functions
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
} 