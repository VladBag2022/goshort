package server

import (
	"fmt"
	"github.com/caarlos0/env/v6"
)

type Config struct {
	Host 	string 	`env:"SERVER_ADDRESS" default:"default-a"`
	Port 	int		`default:"8080"`
	BaseURL string	`env:"BASE_URL"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := env.Parse(&config); err != nil {
		return config, err
	}
	if len(config.BaseURL) == 0 {
		if config.Port != 80 {
			config.BaseURL = fmt.Sprintf("http://%s:%d", config.Host, config.Port)
		} else {
			config.BaseURL = fmt.Sprintf("http://%s", config.Host)
		}
	}
	if config.BaseURL[len(config.BaseURL)-1:] == "/" {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}

	return config, nil
}