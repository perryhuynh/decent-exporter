package config

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"time"
)

type Config struct {
	ListenAddress string
	DecentURL     string
	LogLevel      string
	ReadyMaxAge   time.Duration
	ReconnectMin  time.Duration
	ReconnectMax  time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		ListenAddress: envString("DECENT_EXPORTER_LISTEN_ADDRESS", ":8080"),
		DecentURL:     envString("DECENT_EXPORTER_URL", "http://127.0.0.1:8080"),
		LogLevel:      envString("DECENT_EXPORTER_LOG_LEVEL", "info"),
		ReadyMaxAge:   envDuration("DECENT_EXPORTER_READY_MAX_AGE", 30*time.Second),
		ReconnectMin:  envDuration("DECENT_EXPORTER_RECONNECT_MIN", time.Second),
		ReconnectMax:  envDuration("DECENT_EXPORTER_RECONNECT_MAX", 30*time.Second),
	}
	if _, err := url.ParseRequestURI(cfg.DecentURL); err != nil {
		return Config{}, fmt.Errorf("DECENT_EXPORTER_URL: %w", err)
	}
	if cfg.ReconnectMin <= 0 {
		return Config{}, errors.New("DECENT_EXPORTER_RECONNECT_MIN must be positive")
	}
	if cfg.ReconnectMax < cfg.ReconnectMin {
		return Config{}, errors.New("DECENT_EXPORTER_RECONNECT_MAX must be >= DECENT_EXPORTER_RECONNECT_MIN")
	}
	return cfg, nil
}

func envString(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func envDuration(key string, fallback time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	if parsed, err := time.ParseDuration(value); err == nil {
		return parsed
	}
	if seconds, err := strconv.Atoi(value); err == nil {
		return time.Duration(seconds) * time.Second
	}
	return fallback
}
