// Package config provides configuration management for the application.
package config

import (
	"log"
	"os"
)

type Config struct {
	Port string
}

func getEnv(key string, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}

func Load() Config {
	cfg := Config{
		Port: getEnv("SERVER_PORT", "8080"),
	}

	if cfg.Port == "" {
		log.Fatal("SERVER_PORT is not set")
	}

	return cfg
}
