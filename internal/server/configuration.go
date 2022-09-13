package server

import (
	"fmt"

	"github.com/caarlos0/env/v6"
)

// Config stores application configuration.
type Config struct {
	Address         string `env:"SERVER_ADDRESS" envDefault:"localhost:8080"`
	BaseURL         string `env:"BASE_URL"`
	FileStoragePath string `env:"FILE_STORAGE_PATH"`
	AuthCookieName  string `env:"AUTH_COOKIE" envDefault:"X-AUTH"`
	AuthCookieKey   string `env:"AUTH_KEY" envDefault:"gopher"`
	DatabaseDSN     string `env:"DATABASE_DSN"`
}

// NewConfig parses environment variables and returns config.
func NewConfig() (*Config, error) {
	config := &Config{}
	if err := env.Parse(config); err != nil {
		return nil, err
	}
	if len(config.BaseURL) == 0 {
		config.BaseURL = fmt.Sprintf("http://%s", config.Address)
	}
	if config.BaseURL[len(config.BaseURL)-1:] == "/" {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}
	return config, nil
}
