package configs

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port     string   `yaml:"port"`
	Backends []string `yaml:"backends"`
}

func Load(configPath string) (*Config, error) {
	fullPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("invalid config path: %w", err)
	}

	data, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("config parsing error: %w", err)
	}

	if cfg.Port == "" {
		return nil, fmt.Errorf("port must be specified")
	}

	if len(cfg.Backends) == 0 {
		return nil, fmt.Errorf("at least one backend must be specified")
	}

	return &cfg, nil
}
