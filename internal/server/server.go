// Package server contains HTTP API server.
package server

import (
	"fmt"
	"net/http"

	"github.com/VladBag2022/goshort/internal/storage"
)

type Server struct {
	repository storage.Repository
	postgres   *storage.PostgresRepository
	config     *Config
}

func NewServer(repository storage.Repository, postgres *storage.PostgresRepository, config *Config) Server {
	return Server{
		repository: repository,
		postgres:   postgres,
		config:     config,
	}
}

func (s Server) ListenAndServe() {
	var err error
	if s.config.EnableHTTPS {
		err = http.ListenAndServeTLS(s.config.Address, //nolint:gosec // do not support timeouts for simplicity
			s.config.CertPEMFile,
			s.config.KeyPEMFile,
			router(s))
	} else {
		err = http.ListenAndServe(s.config.Address, router(s)) //nolint:gosec // do not support timeouts for simplicity
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}
