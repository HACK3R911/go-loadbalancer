package env

import (
	"errors"
	"os"
	"strings"
)

const (
	portEnvName     = "PORT"
	backendsEnvName = "BACKENDS"
)

type serverConfig struct {
	port     string
	backends []string
}

func NewServerConfig() (*serverConfig, error) {
	port := os.Getenv(portEnvName)
	if len(port) == 0 {
		return nil, errors.New("PORT not found in environment")
	}

	backendsStr := os.Getenv(backendsEnvName)
	if len(backendsStr) == 0 {
		return nil, errors.New("BACKENDS not found in environment")
	}

	backends := strings.Split(backendsStr, ",")
	for i := range backends {
		backends[i] = strings.TrimSpace(backends[i])
	}

	return &serverConfig{
		port:     port,
		backends: backends,
	}, nil
}

func (cfg *serverConfig) Port() string {
	return cfg.port
}

func (cfg *serverConfig) Backends() []string {
	return cfg.backends
}
