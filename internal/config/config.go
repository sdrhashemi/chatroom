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
		ChatroomUsers   string
	}
	Logging struct {
		Level string
	}
}

func LoadConfig() (*Config, error) {
	// Dynamically locate the .env file in the project root
	projectRoot, err := findProjectRoot()
	if err != nil {
		return nil, err
	}

	envPath := filepath.Join(projectRoot, ".env")
	err = godotenv.Load(envPath)
	if err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	// Populate the config struct from environment variables
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
			ChatroomUsers   string
		}{
			URL:             os.Getenv("NATS_URL"),
			ChatroomSubject: os.Getenv("NATS_CHATROOM_SUBJECT"),
			ChatroomUsers:   os.Getenv("NATS_CHATROOM_USERS"),
		},
		Logging: struct {
			Level string
		}{
			Level: os.Getenv("LOGGING_LEVEL"),
		},
	}

	// Validate required fields
	if cfg.App.ServerURL == "" || cfg.Nats.URL == "" {
		log.Fatalf("Missing required environment variables")
	}

	return cfg, nil
}

// Helper function to locate the project root dynamically
func findProjectRoot() (string, error) {
	// Start from the current working directory
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		// Check if .env exists in the current directory
		if _, err := os.Stat(filepath.Join(dir, ".env")); err == nil {
			return dir, nil
		}

		// Move to the parent directory
		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached the filesystem root
			return "", os.ErrNotExist
		}
		dir = parent
	}
}
