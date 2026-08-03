// Package config loads application configuration from environment
// variables (optionally backed by a local .env file for development).
package config

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Env      string
	HTTPPort string
	LogLevel string

	Database     DatabaseConfig
	Notification NotificationConfig
{% if use_redis %}
	Redis RedisConfig
{% endif %}
{% if use_kafka %}
	Kafka KafkaConfig
{% endif %}
{% if use_s3 %}
	S3 S3Config
{% endif %}
}

type DatabaseConfig struct {
	URL string
}

type NotificationConfig struct {
	BaseURL string
	APIKey  string
}
{% if use_redis %}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}
{% endif %}
{% if use_kafka %}

type KafkaConfig struct {
	Brokers       []string
	ConsumerGroup string
}
{% endif %}
{% if use_s3 %}

type S3Config struct {
	Endpoint  string
	Region    string
	Bucket    string
	AccessKey string
	SecretKey string
}
{% endif %}

// Load reads configuration from the environment. It first loads a .env
// file if one is present (local development convenience) without
// overriding variables already set in the environment.
func Load() (*Config, error) {
	loadDotEnv(".env")

	cfg := &Config{
		Env:      getEnv("APP_ENV", "local"),
		HTTPPort: getEnv("HTTP_PORT", "{{ http_port }}"),
		LogLevel: getEnv("LOG_LEVEL", "info"),
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", ""),
		},
		Notification: NotificationConfig{
			BaseURL: getEnv("NOTIFICATION_BASE_URL", "http://localhost:9999"),
			APIKey:  getEnv("NOTIFICATION_API_KEY", ""),
		},
{% if use_redis %}
		Redis: RedisConfig{
			Addr:     getEnv("REDIS_ADDR", "localhost:6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
{% endif %}
{% if use_kafka %}
		Kafka: KafkaConfig{
			Brokers:       strings.Split(getEnv("KAFKA_BROKERS", "localhost:9092"), ","),
			ConsumerGroup: getEnv("KAFKA_CONSUMER_GROUP", "{{ project_slug }}"),
		},
{% endif %}
{% if use_s3 %}
		S3: S3Config{
			Endpoint:  getEnv("S3_ENDPOINT", "http://localhost:9000"),
			Region:    getEnv("S3_REGION", "us-east-1"),
			Bucket:    getEnv("S3_BUCKET", "{{ project_slug }}"),
			AccessKey: getEnv("S3_ACCESS_KEY", "minioadmin"),
			SecretKey: getEnv("S3_SECRET_KEY", "minioadmin"),
		},
{% endif %}
	}

	if cfg.Database.URL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	var n int
	if _, err := fmt.Sscanf(v, "%d", &n); err != nil {
		return fallback
	}
	return n
}

// loadDotEnv populates os.Environ from a simple KEY=VALUE file, skipping
// blank lines and comments. Existing environment variables always win.
func loadDotEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		key, value, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		key = strings.TrimSpace(key)
		value = strings.Trim(strings.TrimSpace(value), `"'`)
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}
