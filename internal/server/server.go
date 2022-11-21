// Package server contains HTTP API server.
package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/VladBag2022/goshort/internal/storage"
)

type Server struct {
	repository storage.Repository
	postgres   *storage.PostgresRepository
	config     *Config
	http       *http.Server
}

func NewServer(repository storage.Repository, postgres *storage.PostgresRepository, config *Config) Server {
	s := Server{
		repository: repository,
		postgres:   postgres,
		config:     config,
	}
	s.http = &http.Server{ //nolint:gosec // do not support timeouts for simplicity
		Addr:    s.config.Address,
		Handler: router(s),
	}
	return s
}

func (s Server) ListenAndServe() {
	var err error
	if s.config.EnableHTTPS {
		err = s.http.ListenAndServeTLS(s.config.CertPEMFile,
			s.config.KeyPEMFile)
	} else {
		err = s.http.ListenAndServe()
	}
	if err != nil {
		fmt.Println(err)
		return
	}
}

func (s Server) Shutdown() error {
	return s.http.Shutdown(context.Background())
}
