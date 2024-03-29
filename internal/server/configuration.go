package server

import (
	"fmt"

	"github.com/spf13/viper"
)

const (
	DefaultAddress        = "localhost:8080"
	DefaultAuthCookieName = "X-AUTH"
	DefaultAuthCookieKey  = "gopher"
	DefaultCertPEMFile    = "cert.pem"
	DefaultKeyPEMFile     = "key.pem"
	DefaultGRPCAddress    = "localhost:3200"
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
	TrustedSubnet   string
	GRPCAddress     string
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
		TrustedSubnet:   viper.GetString("TrustedSubnet"),
		GRPCAddress:     viper.GetString("GRPCAddress"),
	}
	if len(config.BaseURL) == 0 {
		config.BaseURL = fmt.Sprintf("http://%s", config.Address)
	}
	if config.BaseURL[len(config.BaseURL)-1:] == "/" {
		config.BaseURL = config.BaseURL[:len(config.BaseURL)-1]
	}
	return config
}
