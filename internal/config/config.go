package config

import (
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
)

type Config struct {
	App struct {
		Name      string
		Port      string
		ServerURL string
	}
	Nats struct {
		URL             string
		ChatroomSubject string
	}
}

func LoadConfig() (*Config, error) {
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	cfg := &Config{
		App: struct {
			Name      string
			Port      string
			ServerURL string
		}{
			Name:      os.Getenv("APP_NAME"),
			Port:      os.Getenv("APP_PORT"),
			ServerURL: os.Getenv("APP_SERVER_URL"),
		},
		Nats: struct {
			URL             string
			ChatroomSubject string
		}{
			URL:             os.Getenv("NATS_URL"),
			ChatroomSubject: os.Getenv("NATS_CHATROOM_SUBJECT"),
		},
	}

	// Validate required fields
	if cfg.App.ServerURL == "" || cfg.Nats.URL == "" {
		log.Fatalf("Missing required environment variables")
	}

	return cfg, nil
}

func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
