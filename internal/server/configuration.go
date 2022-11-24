package server

import (
	"fmt"

	"github.com/spf13/viper"
)

// Config stores application configuration.
type Config struct {
	Address         string
	BaseURL         string
	FileStoragePath string
	AuthCookieName  string
	AuthCookieKey   string
	DatabaseDSN     string
	EnableHTTPS     bool
	CertPEMFile     string
	KeyPEMFile      string
	TrustedSubnet      string
}

// NewConfig parses environment variables and returns config.
func NewConfig() *Config {
	config := &Config{
		Address:         viper.GetString("Address"),
		BaseURL:         viper.GetString("BaseURL"),
		FileStoragePath: viper.GetString("FileStoragePath"),
		AuthCookieName:  viper.GetString("AuthCookieName"),
		AuthCookieKey:   viper.GetString("AuthCookieKey"),
		DatabaseDSN:     viper.GetString("DatabaseDSN"),
		EnableHTTPS:     viper.GetBool("EnableHTTPS"),
		CertPEMFile:     viper.GetString("CertPEMFile"),
		TrustedSubnet:      viper.GetString("TrustedSubnet"),
	}
	if len(config.BaseURL) == 0 {
		config.BaseURL = fmt.Sprintf("http://%s", config.Address)
	}
	if config.BaseURL[len(config.BaseURL)-1:] == "/" {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}
	return config
}
