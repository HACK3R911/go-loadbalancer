package yaml

import (
	"errors"

	"github.com/spf13/viper"
)

type serverConfig struct {
	port     string
	backends []string
}

func NewServerConfig(path string) (*serverConfig, error) {
	v := viper.New()

	v.SetConfigFile(path)
	v.SetConfigType("yaml")

	if err := v.ReadInConfig(); err != nil {
		return nil, errors.New("failed to read yaml config file")
	}

	port := v.GetString("port")
	if port == "" {
		return nil, errors.New("port not found in yaml config")
	}

	backends := v.GetStringSlice("backends")
	if len(backends) == 0 {
		return nil, errors.New("backends not found in yaml config")
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
