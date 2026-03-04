package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/HACK3R911/go-loadbalancer/internal/config/env"
	"github.com/HACK3R911/go-loadbalancer/internal/config/yaml"
	"github.com/joho/godotenv"
)

const (
	yamlExt = "yaml"
	envExt  = "env"
)

var (
	ErrPortNotFound      = errors.New("PORT not found")
	ErrBackendsNotFound  = errors.New("BACKENDS not found")
	ErrUnsupportedFormat = errors.New("unsupported config format")
)

func Load(path string) (ServerConfig, error) {
	if path == "" {
		path = "config.yaml"
	}

	format := detectFormat(path)

	switch format {
	case yamlExt:
		return loadYAML(path)
	case envExt:
		return loadENV(path)
	default:
		return nil, ErrUnsupportedFormat
	}
}

func detectFormat(path string) string {
	if path == "" {
		return "yaml"
	}
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".yaml", ".yml":
		return yamlExt
	case ".env":
		return envExt
	default:
		return yamlExt
	}
}

func loadYAML(path string) (ServerConfig, error) {
	return yaml.NewServerConfig(path)
}

func loadENV(path string) (ServerConfig, error) {
	if path == "" {
		path = ".env"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		if os.Getenv("PORT") != "" {
			log.Println("Using environment variables (no .env file found)")
			return env.NewServerConfig()
		}
	}

	if path != ".env" && fileExists(path) {
		if err := godotenv.Load(path); err != nil {
			log.Printf("Ошибка загрузки .env файла: %v", err)
			return nil, err
		}
	}

	return env.NewServerConfig()
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

type ServerConfig interface {
	Port() string
	Backends() []string
}
