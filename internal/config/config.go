package config

import (
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
	configPath := filepath.Join("..", "..", "config.yaml")
	file, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
