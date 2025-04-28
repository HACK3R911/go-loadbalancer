package configs

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Port      string   `yaml:"port"`
	Backends  []string `yaml:"backends"`
	RateLimit struct {
		Capacity   float64 `yaml:"capacity"`
		RefillRate float64 `yaml:"refill_rate"`
	} `yaml:"rate_limit"`
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}
