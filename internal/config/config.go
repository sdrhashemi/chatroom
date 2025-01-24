package config

import (
	"log"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	App struct {
		Name      string `yaml:"name"`
		Port      string `yaml:"port"`
		ServerURL string `yaml:"serverURL"`
	} `yaml:"app"`
	Nats struct {
		URL             string `yaml:"url"`
		ChatroomSubject string `yaml:"chatroomSubject"`
		ChatroomUsers   string `yaml:"chatroomUsers"`
	} `yaml:"nats"`
	Logging struct {
		Level string `yaml:"level"`
	} `yaml:"logging"`
}

func LoadConfig() (*Config, error) {
	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		// Default to a relative path if no environment variable is set
		configPath = "../../config.yaml"
		log.Println("No CONFIG_PATH set, defaulting to", configPath)
	}

	configPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, err
	}

	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}
	log.Println(cfg.App.ServerURL)
	return &cfg, nil
}
